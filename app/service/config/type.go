package config

type IService interface {
	GetRuntimeEnvironment() string
	IsProduction() bool
	GetAPIPort() string
	IsOpenTelemetry() bool

	GetTrelloPropertiesBoardID() string
	GetTrelloInheritanceConfinmentsBoardID() string
	GetTrelloSupportiveDocsBoardID() string
	GetTrelloExpensesBoardID() string

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

	GetInhConfinmentsExcelUpdateWebhook() string
	GetInhConfinmentsNotionUpdateWebhook() string

	GetSupportiveDocsExcelUpdateWebhook() string
	GetSupportiveDocsNotionUpdateWebhook() string
}
