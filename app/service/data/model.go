package data

import (
	"time"

	"github.com/lib/pq"
)

type Property struct {
	ID          int64          `json:"id" db:"id"`
	BoardID     string         `json:"boardId" db:"board_id"`
	CardID      string         `json:"cardId" db:"card_id"`
	Name        string         `json:"name" db:"name"`
	LocationAR  string         `json:"locationAR" db:"location_ar"`
	LocationEN  string         `json:"locationEN" db:"location_en"`
	Lot         string         `json:"lot" db:"lot"`
	Type        string         `json:"type" db:"type"`
	Status      string         `json:"status" db:"status"`
	Owner       string         `json:"owner" db:"owner"`
	Area        float64        `json:"area" db:"area"`
	Shares      float64        `json:"shares" db:"shares"`
	Organized   bool           `json:"organized" db:"is_organized"`
	Effects     bool           `json:"effects" db:"is_effects"`
	Labels      pq.StringArray `json:"labels" db:"labels"`
	Attachments pq.StringArray `json:"attachments" db:"attachments"`
	Comments    pq.StringArray `json:"comments" db:"comments"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

type Attachment struct {
	ID         int64     `json:"id" db:"id"`
	TrelloURL  string    `json:"trelloUrl" db:"trello_url"`
	StorageURL string    `json:"storageUrl" db:"storage_url"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type JobState string

const (
	JobStateQueued    JobState = "queued"
	JobStateRunning   JobState = "running"
	JobStateCancelled JobState = "cancelled"
	JobStateCompleted JobState = "completed"
)

type JobType string

const (
	JobTypeProperties  JobType = "properties"
	JobTypeAttachments JobType = "attachments"
)

type Job struct {
	ID          int64      `json:"id" db:"id"`
	Type        JobType    `json:"type" db:"type"`
	State       JobState   `json:"state" db:"state"`
	Cards       int64      `json:"cards" db:"cards"`
	Errors      int64      `json:"errors" db:"errors"`
	StartedAt   time.Time  `json:"startedAt" db:"started_at"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
}

type Error struct {
	ID         int64     `json:"id" db:"id"`
	Source     string    `json:"source" db:"source"`
	Body       string    `json:"body" db:"body"`
	OccurredAt time.Time `json:"occurredAt" db:"occurred_at"`
}
