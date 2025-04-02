package dto

import (
	"time"

	"survey-project/src/internal/domain"
)

type CreateSurveyInput struct {
	Title       string          `json:"title" validate:"required"`
	Description string          `json:"description"`
	Questions   []QuestionInput `json:"questions" validate:"required,min=1"`
	IsActive    bool            `json:"is_active"`
	IsAnonymous bool            `json:"is_anonymous"`
	RequireInfo bool            `json:"require_info"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
}

type QuestionInput struct {
	Text        string              `json:"text" validate:"required"`
	Type        domain.QuestionType `json:"type" validate:"required,oneof=text number single multiple"`
	Required    bool                `json:"required"`
	Options     []string            `json:"options,omitempty"`
	Min         *int                `json:"min,omitempty"`
	Max         *int                `json:"max,omitempty"`
	Placeholder string              `json:"placeholder,omitempty"`
}

type UpdateSurveyInput struct {
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	Questions   []QuestionInput `json:"questions,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
	IsAnonymous *bool           `json:"is_anonymous,omitempty"`
	RequireInfo *bool           `json:"require_info,omitempty"`
	StartDate   *time.Time      `json:"start_date,omitempty"`
	EndDate     *time.Time      `json:"end_date,omitempty"`
}

type SurveyResponse struct {
	ID          interface{}     `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CreatedBy   string          `json:"created_by"`
	Questions   []QuestionInput `json:"questions"`
	IsActive    bool            `json:"is_active"`
	IsAnonymous bool            `json:"is_anonymous"`
	RequireInfo bool            `json:"require_info"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
