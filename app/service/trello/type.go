package trello

type IService interface {
	RetrieveProperties(max int) ([]TRProperty, error)
	DownloadAttachment(url string) (string, string, string, error)
}
