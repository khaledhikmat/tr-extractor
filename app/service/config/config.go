package config

import (
	"os"
)

type configService struct {
}

func New() IService {
	return &configService{}
}

func (svc *configService) GetRuntimeEnvironment() string {
	if os.Getenv("RUN_TIME_ENV") == "" {
		return "dev"
	}

	return os.Getenv("RUN_TIME_ENV")
}

func (svc *configService) IsProduction() bool {
	return svc.GetRuntimeEnvironment() == "prod"
}

func (svc *configService) GetAPIPort() string {
	return os.Getenv("API_PORT")
}

func (svc *configService) IsOpenTelemetry() bool {
	return os.Getenv("OPEN_TELEMETRY") == "true"
}

func (svc *configService) GetDbDSN() string {
	return os.Getenv("DB_DSN")
}

func (svc *configService) GetTrelloAPIKey() string {
	return os.Getenv("TRELLO_API_KEY")
}

func (svc *configService) GetTrelloToken() string {
	return os.Getenv("TRELLO_TOKEN")
}

func (svc *configService) GetTrelloReadToken() string {
	return os.Getenv("TRELLO_TOKEN_READ")
}

func (svc *configService) GetTrelloBaseURL() string {
	if os.Getenv("TRELLO_BASE_URL") == "" {
		return "https://api.trello.com/1"
	}

	return os.Getenv("TRELLO_BASE_URL")
}

func (svc *configService) GetTrelloDownloadPath() string {
	return os.Getenv("TRELLO_DOWNLOAD_PATH")
}

func (svc *configService) GetTrelloPropertiesBoardID() string {
	if os.Getenv("TRELLO_PROPERTIES_BOARD_ID") == "" {
		return "BXlnOvYt"
	}

	return os.Getenv("TRELLO_PROPERTIES_BOARD_ID")
}

func (svc *configService) GetTrelloInheritanceConfinmentsBoardID() string {
	if os.Getenv("TRELLO_INHERITANCE_CONFINEMENTS_BOARD_ID") == "" {
		return "BXlnOvYt"
	}

	return os.Getenv("TRELLO_INHERITANCE_CONFINEMENTS_BOARD_ID")
}

func (svc *configService) GetTrelloSupportiveDocsBoardID() string {
	if os.Getenv("TRELLO_SUPPORTIVE_DOCS_BOARD_ID") == "" {
		return "bOEmEE4S"
	}

	return os.Getenv("TRELLO_SUPPORTIVE_DOCS_BOARD_ID")
}

func (svc *configService) GetTrelloExpensesBoardID() string {
	if os.Getenv("TRELLO_EXPENSES_BOARD_ID") == "" {
		return "bOEmEE4S"
	}

	return os.Getenv("TRELLO_EXPENSES_BOARD_ID")
}

func (svc *configService) GetDropboxAccessToken() string {
	return os.Getenv("DROPBOX_ACCESS_TOKEN")
}

func (svc *configService) GetDropboxUploadPath() string {
	return os.Getenv("DROPBOX_UPLOAD_PATH")
}

func (svc *configService) GetStorageBucket() string {
	return os.Getenv("STORAGE_BUCKET")
}

func (svc *configService) GetStorageRegion() string {
	return os.Getenv("STORAGE_REGION")
}

func (svc *configService) GetPropertiesExcelUpdateWebhook() string {
	if os.Getenv("PROPERTIES_EXCEL_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/bk7ct9twnq3sndfj4kmc7sx2idhjbukf"
	}

	return os.Getenv("PROPERTIES_EXCEL_UPDATE_WEBHOOK")
}

func (svc *configService) GetPropertiesNotionUpdateWebhook() string {
	if os.Getenv("PROPERTIES_NOTION_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/sq2jd7e2zjn2hobklhkigln4ivrb73il"
	}

	return os.Getenv("PROPERTIES_NOTION_UPDATE_WEBHOOK")
}

func (svc *configService) GetInhConfinmentsExcelUpdateWebhook() string {
	if os.Getenv("INH_CONFINMENTS_EXCEL_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/add"
	}

	return os.Getenv("INH_CONFINMENTS_EXCEL_UPDATE_WEBHOOK")
}

func (svc *configService) GetInhConfinmentsNotionUpdateWebhook() string {
	if os.Getenv("INH_CONFINMENTS_NOTION_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/add"
	}

	return os.Getenv("INH_CONFINMENTS_NOTION_UPDATE_WEBHOOK")
}

func (svc *configService) GetSupportiveDocsExcelUpdateWebhook() string {
	if os.Getenv("SUPPORTIVE_DOCS_EXCEL_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/add"
	}

	return os.Getenv("SUPPORTIVE_DOCS_EXCEL_UPDATE_WEBHOOK")
}

func (svc *configService) GetSupportiveDocsNotionUpdateWebhook() string {
	if os.Getenv("SUPPORTIVE_DOCS_NOTION_UPDATE_WEBHOOK") == "" {
		return "https://hook.us2.make.com/add"
	}

	return os.Getenv("SUPPORTIVE_DOCS_NOTION_UPDATE_WEBHOOK")
}
