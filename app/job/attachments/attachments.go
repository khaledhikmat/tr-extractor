package jobattachments

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/khaledhikmat/tr-extractor/service/config"
	"github.com/khaledhikmat/tr-extractor/service/data"
	"github.com/khaledhikmat/tr-extractor/service/lgr"
	"github.com/khaledhikmat/tr-extractor/service/storage"
	"github.com/khaledhikmat/tr-extractor/service/trello"
)

func Processor(ctx context.Context,
	jobID int64,
	pageSize int,
	errorStream chan error,
	cfgsvc config.IService,
	datasvc data.IService,
	trlsvc trello.IService,
	storagesvc storage.IService) {

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
	attachments := []string{}
	finalState := data.JobStateCompleted

	defer func() {
		// Update job state to completed
		now := time.Now()
		job.State = finalState
		job.Cards = int64(len(attachments))
		job.Errors = int64(errors)
		job.CompletedAt = &now
		err = datasvc.UpdateJob(&job)
		if err != nil {
			errorStream <- err
			return
		}
	}()

	// Retrieve property attachments from the database
	propatts, err := datasvc.RetrievePropertyAttachments(pageSize)
	if err != nil {
		errorStream <- err
		return
	}
	attachments = append(attachments, propatts...)

	// TODO: Add other attachments from Trello
	//attachments = append(attachments, inhatts...)

	// Insert/update property attachments into the database
	for _, attachmentURL := range attachments {

		// TODO: This should not happen!!!
		if attachmentURL == "" {
			continue
		}

		parts := strings.SplitN(attachmentURL, "|", 3)
		if len(parts) != 3 {
			errorStream <- fmt.Errorf("%s. URL is not constructed properly. No '|' found in the input string", attachmentURL)
			errors++
			continue
		}

		// In order to add some context to URLs, I appended property name to the URL
		// after the pipe character
		attachmentURL = parts[0]
		attachmentFolder := parts[1]
		attachmentPostfix := parts[2]

		// If the context is cancelled, exit the loop
		// But execute the defer block first
		select {
		case <-ctx.Done():
			finalState = data.JobStateCancelled
			return
		default:
		}

		// If the attachment is already uploaded, skip it
		uploaded, err := datasvc.IsAttachmentMapped(attachmentURL)
		if err != nil {
			errorStream <- err
			errors++
			continue
		}

		if uploaded {
			continue
		}

		// Download the attachment from Trello
		localPath, attachmentID, extension, err := trlsvc.DownloadAttachment(attachmentURL)
		if err != nil {
			errorStream <- err
			errors++
			continue
		}

		// Upload to Cloud Storage
		cloudURL, err := storagesvc.Upload(localPath, attachmentFolder, fmt.Sprintf("%s-%s%s", attachmentID, attachmentPostfix, extension))
		if err != nil {
			errorStream <- err
			errors++
			continue
		}

		// Map attachment to Cloud URL
		err = datasvc.MapAttachment(attachmentURL, cloudURL)
		if err != nil {
			errorStream <- err
			errors++
			continue
		}
	}

	lgr.Logger.Debug("jobattachments.Processor",
		slog.String("event", "done"),
	)
}
