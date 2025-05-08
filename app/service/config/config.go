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

func (svc *configService) GetTrelloBaseURL() string {
	if os.Getenv("TRELLO_BASE_URL") == "" {
		return "https://api.trello.com/1"
	}

	return os.Getenv("TRELLO_BASE_URL")
}

func (svc *configService) GetTrelloPropertiesBoardID() string {
	if os.Getenv("TRELLO_PROPERTIES_BOARD_ID") == "" {
		return "5f3c2b4e1a0d3b2f8c4e4d6f"
	}

	return os.Getenv("TRELLO_PROPERTIES_BOARD_ID")
}

func (svc *configService) GetTrelloExpensesBoardID() string {
	if os.Getenv("TRELLO_EXPENSES_BOARD_ID") == "" {
		return "5f3c2b4e1a0d3b2f8c4e4d6f"
	}

	return os.Getenv("TRELLO_EXPENSES_BOARD_ID")
}
func (svc *configService) GetTrelloInheritanceConfinementsBoardID() string {
	if os.Getenv("TRELLO_INHERITANCE_CONFINEMENTS_BOARD_ID") == "" {
		return "5f3c2b4e1a0d3b2f8c4e4d6f"
	}

	return os.Getenv("TRELLO_INHERITANCE_CONFINEMENTS_BOARD_ID")
}
func (svc *configService) GetTrelloToDosBoardID() string {
	if os.Getenv("TRELLO_TODO_BOARD_ID") == "" {
		return "5f3c2b4e1a0d3b2f8c4e4d6f"
	}

	return os.Getenv("TRELLO_TODO_BOARD_ID")
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
