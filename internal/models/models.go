package models

import "time"

type Good struct {
	ID          int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type DeletedGood struct {
	ID        int  `json:"id" db:"id"`
	ProjectID int  `json:"projectId" db:"campaign_id"`
	Removed   bool `json:"removed" db:"removed"`
}

type Priorities struct {
	ID       int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}

type GoodPriorities struct {
	Priorities []Priorities `json:"priorities"`
}
