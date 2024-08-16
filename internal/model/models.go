package model

import (
	"time"
)

// CloudEvent represents the structure of a CloudEvent.
type CloudEvent struct {
	SpecVersion     string      `json:"specversion"`
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            time.Time   `json:"time"`
	Data            interface{} `json:"data"`
	ObjectKey       string      `json:"object_key"` // 파일 이름을 추가
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

type MessageData struct {
	AuthorID    int64     `json:"author_id"`
	BoardID     int64     `json:"board_id"`
	Content     string    `json:"content"`
	CreatedOn   time.Time `json:"created_on"`
	ID          int64     `json:"id"`
	LastReplyID ID        `json:"last_reply_id"`
	Locked      bool      `json:"locked"`
	ParentID    ID        `json:"parent_id"`
	Sticky      bool      `json:"sticky"`
	Subject     string    `json:"subject"`
	UpdatedOn   time.Time `json:"updated_on"`
}

type ID struct {
	Int64 int64 `json:"Int64"`
	Valid bool  `json:"Valid"`
}

type UserData struct {
	Admin            bool         `json:"admin"`
	AuthSourceID     AuthSourceID `json:"auth_source_id"`
	CreatedOn        time.Time    `json:"created_on"`
	Firstname        string       `json:"firstname"`
	HashedPassword   string       `json:"hashed_password"`
	ID               int64        `json:"id"`
	Language         string       `json:"language"`
	LastLoginOn      On           `json:"last_login_on"`
	Lastname         string       `json:"lastname"`
	Login            string       `json:"login"`
	MailNotification string       `json:"mail_notification"`
	MustChangePasswd bool         `json:"must_change_passwd"`
	PasswdChangedOn  On           `json:"passwd_changed_on"`
	Salt             string       `json:"salt"`
	Status           int64        `json:"status"`
	Type             string       `json:"type"`
	UpdatedOn        time.Time    `json:"updated_on"`
}

type AuthSourceID struct {
	Int64 int64 `json:"Int64"`
	Valid bool  `json:"Valid"`
}

type On struct {
	Time  time.Time `json:"Time"`
	Valid bool      `json:"Valid"`
}

type IssueData struct {
	ID             int       `json:"id"`
	JobID          int       `json:"job_id"`
	TrackerID      int       `json:"tracker_id"`
	ProjectID      int       `json:"project_id"`
	Subject        string    `json:"subject"`
	Description    string    `json:"description"`
	DueDate        time.Time `json:"due_date"`
	StatusID       int       `json:"status_id"`
	AssignedToID   int       `json:"assigned_to_id"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	StartDate      time.Time `json:"start_date"`
	DoneRatio      int       `json:"done_ratio"`
	EstimatedHours float64   `json:"estimated_hours"`
	PriorityID     int       `json:"priority_id"`
	AuthorID       int       `json:"author_id"`
	CommentorID    int       `json:"commentor_id"`
	RootID         int       `json:"root_id"`
	Notes          string    `json:"notes"`
	Property       string    `json:"property"`
	PropKey        string    `json:"prop_key"`
	OldValue       string    `json:"old_value"`
	Value          string    `json:"value"`
}
