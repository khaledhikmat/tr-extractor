package jobproperties

import (
	"context"
	"log/slog"
	"time"

	jobb "github.com/khaledhikmat/tr-extractor/job"
	"github.com/khaledhikmat/tr-extractor/service/config"
	"github.com/khaledhikmat/tr-extractor/service/data"
	"github.com/khaledhikmat/tr-extractor/service/lgr"
	"github.com/khaledhikmat/tr-extractor/service/storage"
	"github.com/khaledhikmat/tr-extractor/service/trello"
	"github.com/khaledhikmat/tr-extractor/utils"
)

func Processor(ctx context.Context,
	jobID int64,
	pageSize int,
	errorStream chan error,
	cfgsvc config.IService,
	datasvc data.IService,
	trsvc trello.IService,
	_ storage.IService) {

	// Update job state to running
	job, err := datasvc.RetrieveJobByID(jobID)
	if err != nil {
		errorStream <- err
		return
	}
	job.State = data.JobStateRunning
	err = datasvc.UpdateJob(&job)
	if err != nil {
		errorStream <- err
		return
	}

	errors := 0
	trprops := []trello.TRProperty{}
	finalState := data.JobStateCompleted
	boardID := cfgsvc.GetTrelloPropertiesBoardID()

	defer func() {
		// Update job state to completed
		now := time.Now()
		job.State = finalState
		job.Cards = int64(len(trprops))
		job.Errors = int64(errors)
		job.CompletedAt = &now
		err = datasvc.UpdateJob(&job)
		if err != nil {
			errorStream <- err
			return
		}
	}()

	// Retrieve properties from Trello
	trprops, err = trsvc.RetrieveProperties(pageSize)
	if err != nil {
		errorStream <- err
		errors++
	}

	// Insert/update properties into the database
	for _, trprop := range trprops {

		// If the context is cancelled, exit the loop
		// But execute the defer block first
		select {
		case <-ctx.Done():
			finalState = data.JobStateCancelled
			return
		default:
		}

		// If Trello's last activity date is zero, use the current time
		updatedAt := time.Now()
		if !trprop.DateLastActivity.IsZero() {
			updatedAt = trprop.DateLastActivity
		}

		// Convert to data model property
		prop := data.Property{
			BoardID:    boardID,
			CardID:     trprop.ID,
			Name:       trprop.Name,
			LocationAR: trprop.LocationAR,
			LocationEN: trprop.LocationEN,
			Lot:        trprop.Lot,
			Type:       trprop.Type,
			Status:     trprop.Status,
			Owner:      trprop.Owner,
			Area:       trprop.Area,
			Shares:     trprop.Shares,
			Organized:  trprop.Organized,
			Effects:    trprop.Effects,
			Labels: utils.Map(trprop.Labels, func(label trello.TRLabel) string {
				return label.Name
			}),
			Attachments: utils.Map(trprop.Attachments, func(attachment trello.TRAttachment) string {
				// return fmt.Sprintf("%s|%s|%s|%s", trprop.ID, attachment.ID, trprop.Name, attachment.URL)
				return attachment.URL
			}),
			Comments: utils.Map(trprop.Comments, func(comment trello.TRComment) string {
				return comment.Data.Text
			}),
			UpdatedAt: updatedAt,
		}

		// Insert or update the property into the database
		_, _, err := datasvc.NewProperty(prop)
		if err != nil {
			errorStream <- err
			errors++
			continue
		}
	}

	lgr.Logger.Debug("jobproperties.Processor",
		slog.String("event", "done"),
	)

	// Notify the automation webhook to trigger
	// lgr.Logger.Debug("jobproperties.Processor",
	// 	slog.String("webhookUrl", cfgsvc.GetPropertiesExcelUpdateWebhook()),
	// )
	// err = jobb.PostToAutomationWebhook(cfgsvc.GetPropertiesExcelUpdateWebhook())
	// if err != nil {
	// 	errorStream <- err
	// }

	// Notify the automation webhook to trigger
	lgr.Logger.Debug("jobproperties.Processor",
		slog.String("webhookUrl", cfgsvc.GetPropertiesNotionUpdateWebhook()),
	)
	err = jobb.PostToAutomationWebhook(cfgsvc.GetPropertiesNotionUpdateWebhook())
	if err != nil {
		errorStream <- err
	}
}
