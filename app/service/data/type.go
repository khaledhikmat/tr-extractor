package data

type IService interface {
	ResetFactory() error

	NewProperty(prop Property) (bool, int64, error)
	UpdateProperty(prop *Property) error
	RetrieveProperties(page, pageSize int, orderBy, orderDir string) ([]Property, error)
	RetrievePropertyAttachments(pageSize int) ([]string, error)

	IsAttachmentMapped(url string) (bool, error)
	MapAttachment(trelloURL, storageURL string) error

	NewJob(job Job) (int64, error)
	UpdateJob(job *Job) error
	RetrieveJobByID(id int64) (Job, error)
	IsPendingJobsByType(jobType JobType) (bool, error)

	NewAPIKey(key string) error
	IsAPIKeyValid(key string) (bool, error)
	NewError(source, body string) error
}
