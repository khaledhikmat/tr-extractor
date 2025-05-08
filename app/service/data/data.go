package data

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Import the PostgreSQL driver

	"github.com/khaledhikmat/tr-extractor/service/config"
)

var mutex = &sync.Mutex{}

//go:embed sql/reset_factory.sql
var resetfactorySQL string

//go:embed sql/insertproperty.sql
var insertpropertySQL string

//go:embed sql/updateproperty.sql
var updatepropertySQL string

//go:embed sql/insertjob.sql
var insertjobSQL string

//go:embed sql/updatejob.sql
var updatejobSQL string

//go:embed sql/insertapikey.sql
var insertapikeySQL string

//go:embed sql/inserterror.sql
var inserterrorSQL string

type dataService struct {
	ConfigSvc config.IService
	Db        *sqlx.DB
}

func New(cfgsvc config.IService) IService {
	return &dataService{
		ConfigSvc: cfgsvc,
	}
}

func (svc *dataService) ResetFactory() error {
	err := svc.dbConnection()
	if err != nil {
		return err
	}

	_, err = svc.Db.Exec(resetfactorySQL)
	if err != nil {
		return err
	}

	return nil
}

func (svc *dataService) NewProperty(prop Property) (bool, int64, error) {
	err := svc.dbConnection()
	if err != nil {
		return false, -1, err
	}

	// Make sure that property does not already exist
	p, err := svc.retrievePropertyByIDs(prop.BoardID, prop.CardID)
	if err != nil {
		return false, -1, fmt.Errorf("Error fetching property by ID: %v", err)
	}

	// If the property already exists, switch to update the attributes
	if p.CardID != "" {
		err = svc.UpdateProperty(&prop)
		return false, prop.ID, err
	}

	// Convert to args so it can be used with the database
	args := map[string]interface{}{
		"board_id":     prop.BoardID,
		"card_id":      prop.CardID,
		"name":         prop.Name,
		"location_ar":  prop.LocationAR,
		"location_en":  prop.LocationEN,
		"lot":          prop.Lot,
		"type":         prop.Type,
		"status":       prop.Status,
		"owner":        prop.Owner,
		"area":         prop.Area,
		"shares":       prop.Shares,
		"is_organized": prop.Organized,
		"is_effects":   prop.Effects,
		"labels":       pq.Array(prop.Labels),
		"attachments":  pq.Array(prop.Attachments),
		"comments":     pq.Array(prop.Comments),
	}

	// Execute the insert query using NamedExec or NamedQuery
	rows, err := svc.Db.NamedQuery(insertpropertySQL, args)
	if err != nil {
		return false, -1, err
	}
	defer rows.Close()

	// Fetch the newly inserted ID if needed
	if rows.Next() {
		err = rows.Scan(&prop.ID)
		if err != nil {
			return false, -1, err
		}
	}

	return true, prop.ID, nil
}

func (svc *dataService) UpdateProperty(prop *Property) error {
	err := svc.dbConnection()
	if err != nil {
		return err
	}

	// Make sure the property does exist
	p, err := svc.retrievePropertyByIDs(prop.BoardID, prop.CardID)
	if err != nil {
		return fmt.Errorf("error fetching property by ID: %v", err)
	}

	if p.CardID == "" {
		return fmt.Errorf("card ID %s does not exist", p.CardID)
	}

	_, err = svc.Db.Exec(
		updatepropertySQL,
		prop.BoardID,
		prop.CardID,
		prop.Name,
		prop.LocationAR,
		prop.LocationEN,
		prop.Lot,
		prop.Type,
		prop.Status,
		prop.Owner,
		prop.Area,
		prop.Shares,
		prop.Organized,
		prop.Effects,
		pq.Array(prop.Labels),
		pq.Array(prop.Attachments),
		pq.Array(prop.Comments),
		p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *dataService) RetrieveProperties(page, pageSize int, orderBy, orderDir string) ([]Property, error) {
	props := []Property{}
	err := svc.dbConnection()
	if err != nil {
		return props, err
	}

	if page < 1 {
		return props, fmt.Errorf("Invalid page number %d", page)
	}

	if pageSize <= 0 {
		return props, fmt.Errorf("Invalid page size %d", pageSize)
	}

	if orderBy != "updated_at" &&
		orderBy != "area" &&
		orderBy != "comments" &&
		orderBy != "attachments" {
		return props, fmt.Errorf("Invalid order by %s", orderBy)
	}

	if orderDir != "asc" && orderDir != "desc" {
		return props, fmt.Errorf("Invalid order direction %s", orderDir)
	}

	boardID := svc.ConfigSvc.GetTrelloPropertiesBoardID()

	// Calculate the offset
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
        SELECT * 
		FROM properties 
		WHERE board_id = $1 
		ORDER BY %s %s 
		LIMIT $2 OFFSET $3 
    `, orderBy, orderDir)

	err = svc.Db.Select(&props, query, boardID, pageSize, offset)
	if err != nil {
		return props, err
	}

	return props, nil
}

func (svc *dataService) retrievePropertyByIDs(boardID string, propID string) (Property, error) {
	err := svc.dbConnection()
	if err != nil {
		return Property{}, err
	}

	var props []Property
	query := `
        SELECT * 
		FROM properties
		WHERE board_id = $1 
		AND card_id = $2 
		LIMIT 1
    `

	err = svc.Db.Select(&props, query, boardID, propID)
	if err != nil {
		return Property{}, err
	}

	if len(props) == 0 {
		return Property{}, nil
	}

	return props[0], nil
}

func (svc *dataService) NewJob(job Job) (int64, error) {
	err := svc.dbConnection()
	if err != nil {
		return -1, err
	}

	// Execute the insert query using NamedExec or NamedQuery
	rows, err := svc.Db.NamedQuery(insertjobSQL, job)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	// Fetch the newly inserted ID if needed
	if rows.Next() {
		err = rows.Scan(&job.ID)
		if err != nil {
			return -1, err
		}
	}

	return job.ID, nil
}

func (svc *dataService) UpdateJob(job *Job) error {
	err := svc.dbConnection()
	if err != nil {
		return err
	}

	_, err = svc.Db.Exec(updatejobSQL, job.State, job.Cards, job.Errors, job.CompletedAt, job.ID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *dataService) RetrieveJobByID(id int64) (Job, error) {
	err := svc.dbConnection()
	if err != nil {
		return Job{}, err
	}

	var jobs []Job
	query := `
        SELECT * FROM jobs 
		WHERE id = $1 
		LIMIT 1
    `

	err = svc.Db.Select(&jobs, query, id)
	if err != nil {
		return Job{}, err
	}

	if len(jobs) == 0 {
		return Job{}, fmt.Errorf("Job ID %d does not exist", id)
	}

	return jobs[0], nil
}

func (svc *dataService) IsPendingJobsByType(jobType JobType) (bool, error) {
	err := svc.dbConnection()
	if err != nil {
		return false, err
	}

	var jobs []Job
	query := `
        SELECT * FROM jobs 
		WHERE  type = $1
		AND state IN ($2, $3) 
		LIMIT 1
    `

	err = svc.Db.Select(&jobs, query, jobType, JobStateQueued, JobStateRunning)
	if err != nil {
		return false, err
	}

	return len(jobs) > 0, nil
}

func (svc *dataService) NewAPIKey(key string) error {
	err := svc.dbConnection()
	if err != nil {
		return err
	}

	// Execute the insert query using NamedExec or NamedQuery
	_, err = svc.Db.Exec(insertapikeySQL, key)
	if err != nil {
		return err
	}

	return nil
}

func (svc *dataService) IsAPIKeyValid(key string) (bool, error) {
	err := svc.dbConnection()
	if err != nil {
		return false, err
	}

	var keys []string
	query := `
        SELECT key 
		FROM api_keys 
		WHERE key = $1
		AND expires_at > now() 
		LIMIT 1
    `

	err = svc.Db.Select(&keys, query, key)
	if err != nil {
		return false, err
	}

	if len(keys) == 0 {
		return false, fmt.Errorf("APP KEY %s is not valid", key)
	}

	return true, nil
}

func (svc *dataService) NewError(source, body string) error {
	err := svc.dbConnection()
	if err != nil {
		return err
	}

	// Execute the insert query using NamedExec or NamedQuery
	_, err = svc.Db.Exec(inserterrorSQL, source, body)
	if err != nil {
		return err
	}

	return nil
}

func (svc *dataService) Finalize() {
	if svc.Db != nil {
		svc.Db.Close()
	}
}

func (svc *dataService) dbConnection() error {
	var err error
	if svc.Db != nil {
		return nil
	}

	// Allow one thread to access the database at a time
	mutex.Lock()
	defer mutex.Unlock()

	svc.Db, err = sqlx.Connect("postgres", svc.ConfigSvc.GetDbDSN())
	if err != nil {
		return err
	}

	return nil
}
