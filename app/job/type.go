package job

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/khaledhikmat/tr-extractor/service/config"
	"github.com/khaledhikmat/tr-extractor/service/data"
	"github.com/khaledhikmat/tr-extractor/service/storage"
	"github.com/khaledhikmat/tr-extractor/service/trello"
)

// Signature of job processors
type Processor func(ctx context.Context,
	jobId int64,
	pageSize int,
	errorStream chan error,
	cfgsvc config.IService,
	datasvc data.IService,
	trsvc trello.IService,
	storagesvc storage.IService)

func PostToAutomationWebhook(url string) error {
	if url == "" {
		return fmt.Errorf("postToAutomationWebhook - automation webhook URL is empty")
	}

	// Provide a dummy payload
	payload := ""
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("postToAutomationWebhook - could not create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "text/plain")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("postToAutomationWebhook - request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("postToAutomationWebhook - unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
