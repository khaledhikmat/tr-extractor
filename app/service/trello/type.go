package trello

type IService interface {
	RetrieveProperties(max int) ([]TRProperty, error)
}
