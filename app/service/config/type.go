package config

type IService interface {
	GetRuntimeEnvironment() string
	IsProduction() bool
	GetAPIPort() string
	IsOpenTelemetry() bool

	GetTrelloPropertiesBoardID() string
	GetTrelloExpensesBoardID() string
	GetTrelloInheritanceConfinementsBoardID() string
	GetTrelloToDosBoardID() string

	GetDbDSN() string
	GetTrelloAPIKey() string
	GetTrelloToken() string
	GetTrelloBaseURL() string

	GetPropertiesExcelUpdateWebhook() string
	GetPropertiesNotionUpdateWebhook() string
}
