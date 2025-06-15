package data

type IService interface {
	ResetFactory() error

	NewProperty(prop Property) (bool, int64, error)
	UpdateProperty(prop *Property) error
	RetrieveProperties(page, pageSize int, orderBy, orderDir string) ([]Property, error)
	RetrievePropertyAttachments(pageSize int) ([]string, error)

	NewInheritanceConfinment(inh InheritanceConfinment) (bool, int64, error)
	UpdateInheritanceConfinment(inh *InheritanceConfinment) error
	RetrieveInheritanceConfinments(page, pageSize int, orderBy, orderDir string) ([]InheritanceConfinment, error)
	RetrieveInheritanceConfinmentAttachments(pageSize int) ([]string, error)

	NewSupportiveDoc(inh SupportiveDoc) (bool, int64, error)
	UpdateSupportiveDoc(inh *SupportiveDoc) error
	RetrieveSupportiveDocs(page, pageSize int, orderBy, orderDir string) ([]SupportiveDoc, error)
	RetrieveSupportiveDocAttachments(pageSize int) ([]string, error)

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
