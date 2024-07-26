package model

import "time"

// CloudEvent represents the structure of a CloudEvent.
type CloudEvent struct {
	SpecVersion     string    `json:"specversion"`
	ID              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            EventData `json:"data"`
	ObjectKey       string    `json:"object_key"` // 파일 이름을 추가
}

// EventData represents the data contained within a CloudEvent.
type EventData struct {
	ID             int       `json:"id"`
	JobID          int       `json:"job_id"`
	Status         string    `json:"status"`
	Assignee       string    `json:"assignee"`
	StartDate      time.Time `json:"start_date"`
	DueDate        time.Time `json:"due_date"`
	DoneRatio      int       `json:"done_ratio"`
	EstimatedHours int       `json:"estimated_hours"`
	Priority       string    `json:"priority"`
	Author         string    `json:"author"`
	Subject        string    `json:"subject"`
	Description    string    `json:"description"`
	Commentor      string    `json:"commentor"`
	Notes          string    `json:"notes"`
	CreatedOn      time.Time `json:"created_on"`
}
