package domain

import (
	"time"
)

type QuestionType string

const (
	QuestionTypeText     QuestionType = "text"
	QuestionTypeNumber   QuestionType = "number"
	QuestionTypeSingle   QuestionType = "single"
	QuestionTypeMultiple QuestionType = "multiple"
)

type Question struct {
	ID          interface{}  `bson:"id" json:"id"`
	Text        string       `bson:"text" json:"text"`
	Type        QuestionType `bson:"type" json:"type"`
	Required    bool         `bson:"required" json:"required"`
	Options     []string     `bson:"options,omitempty" json:"options,omitempty"`
	Min         *int         `bson:"min,omitempty" json:"min,omitempty"`
	Max         *int         `bson:"max,omitempty" json:"max,omitempty"`
	Placeholder string       `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
}

type RespondentInfo struct {
	FullName string `bson:"full_name,omitempty" json:"full_name,omitempty"`
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
}

type Survey struct {
	ID          interface{} `bson:"id,omitempty" json:"id"`
	Title       string      `bson:"title" json:"title"`
	Description string      `bson:"description" json:"description"`
	CreatedBy   string      `bson:"created_by" json:"created_by"`
	Questions   []Question  `bson:"questions" json:"questions"`
	IsActive    bool        `bson:"is_active" json:"is_active"`
	IsAnonymous bool        `bson:"is_anonymous" json:"is_anonymous"`
	RequireInfo bool        `bson:"require_info" json:"require_info"`
	StartDate   time.Time   `bson:"start_date" json:"start_date"`
	EndDate     time.Time   `bson:"end_date" json:"end_date"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}

type SurveyRepository interface {
	Create(survey *Survey) error
	GetByID(id string) (*Survey, error)
	GetAll() ([]*Survey, error)
	Update(survey *Survey) error
	Delete(id string) error
	GetByUserID(userID string) ([]*Survey, error)
}
