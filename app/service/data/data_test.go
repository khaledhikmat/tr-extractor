package data

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/khaledhikmat/tr-extractor/service/config"
)

func TestDataRetrieveAttachments(t *testing.T) {

	// TODO: load a fake configuration service instead of loading the .env file
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
		return
	}

	configSvc := config.New()
	dataSvc := New(configSvc)

	urls, err := dataSvc.RetrievePropertyAttachments(10)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Retrieved attachments:", len(urls))
	if len(urls) < 20 {
		t.Error(fmt.Errorf("expected less than urls, got %d", len(urls)))
		return
	}
}
