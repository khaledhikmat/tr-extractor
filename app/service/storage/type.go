package storage

type IService interface {
	Upload(filePath, folder, indentifier string) (string, error)
}
