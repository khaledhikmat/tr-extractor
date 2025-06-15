package jobsupportivedocs

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
	trprops := []trello.TRSupportiveDoc{}
	finalState := data.JobStateCompleted
	boardID := cfgsvc.GetTrelloSupportiveDocsBoardID()

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

	// Retrieve supportive docs from Trello
	trprops, err = trsvc.RetrieveSupportiveDocs(pageSize)
	if err != nil {
		errorStream <- err
		errors++
	}

	// Insert/update inhconfs into the database
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

		if trprop.Title == "" {
			trprop.Title = trprop.Name
		}

		// Convert to data model inhconf
		prop := data.SupportiveDoc{
			BoardID:  boardID,
			CardID:   trprop.ID,
			Name:     trprop.Name,
			Title:    trprop.Title,
			Category: trprop.Category,
			Labels: utils.Map(trprop.Labels, func(label trello.TRLabel) string {
				return label.Name
			}),
			Attachments: utils.Map(trprop.Attachments, func(attachment trello.TRAttachment) string {
				return attachment.URL
			}),
			Comments: utils.Map(trprop.Comments, func(comment trello.TRComment) string {
				return comment.Data.Text
			}),
			UpdatedAt: updatedAt,
		}

		// Insert or update the supportive doc into the database
		_, _, err := datasvc.NewSupportiveDoc(prop)
		if err != nil {
			errorStream <- err
			errors++
			continue
		}
	}

	lgr.Logger.Debug("jobsupportivedocs.Processor",
		slog.String("event", "done"),
	)

	// Notify the automation webhook to trigger
	// lgr.Logger.Debug("supportivedocsconfs.Processor",
	// 	slog.String("webhookUrl", cfgsvc.GetSupportiveDocsExcelUpdateWebhook()),
	// )
	// err = jobb.PostToAutomationWebhook(cfgsvc.GetSupportiveDocsExcelUpdateWebhook())
	// if err != nil {
	// 	errorStream <- err
	// }

	// Notify the automation webhook to trigger
	lgr.Logger.Debug("jobsupportivedocs.Processor",
		slog.String("webhookUrl", cfgsvc.GetSupportiveDocsNotionUpdateWebhook()),
	)
	err = jobb.PostToAutomationWebhook(cfgsvc.GetSupportiveDocsNotionUpdateWebhook())
	if err != nil {
		errorStream <- err
	}
}
