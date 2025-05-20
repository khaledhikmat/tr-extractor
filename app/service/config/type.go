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
	GetTrelloReadToken() string
	GetTrelloBaseURL() string
	GetTrelloDownloadPath() string

	GetDropboxAccessToken() string
	GetDropboxUploadPath() string

	GetStorageBucket() string
	GetStorageRegion() string

	GetPropertiesExcelUpdateWebhook() string
	GetPropertiesNotionUpdateWebhook() string
}
