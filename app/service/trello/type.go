package trello

type IService interface {
	RetrieveProperties(max int) ([]TRProperty, error)
	RetrieveInheritanceConfinments(max int) ([]TRInheritanceConfinement, error)
	RetrieveSupportiveDocs(max int) ([]TRSupportiveDoc, error)
	DownloadAttachment(url string) (string, string, string, error)
}
