package usecase

import (
	"reflect"
	"time"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/domain"
	"survey-project/src/internal/dto"
)

type SurveyUsecase struct {
	surveyRepo domain.SurveyRepository
}

func NewSurveyUsecase(surveyRepo domain.SurveyRepository) *SurveyUsecase {
	return &SurveyUsecase{
		surveyRepo: surveyRepo,
	}
}

func (u *SurveyUsecase) Create(input dto.CreateSurveyInput, userID string) (*dto.SurveyResponse, error) {
	questions := make([]domain.Question, len(input.Questions))
	for i, q := range input.Questions {
		questions[i] = domain.Question{
			Text:        q.Text,
			Type:        q.Type,
			Required:    q.Required,
			Options:     q.Options,
			Min:         q.Min,
			Max:         q.Max,
			Placeholder: q.Placeholder,
		}
	}

	survey := &domain.Survey{
		Title:       input.Title,
		Description: input.Description,
		CreatedBy:   userID,
		Questions:   questions,
		IsActive:    input.IsActive,
		IsAnonymous: input.IsAnonymous,
		RequireInfo: input.RequireInfo,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.surveyRepo.Create(survey); err != nil {
		return nil, err
	}

	return convertToResponse(survey), nil
}

func (u *SurveyUsecase) GetByID(id string, userID string, role string) (*dto.SurveyResponse, error) {
	survey, err := u.surveyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if role != string(domain.RoleAdmin) && survey.CreatedBy != userID {
		return nil, apperrors.ErrUnauthorized
	}

	return convertToResponse(survey), nil
}

func (u *SurveyUsecase) GetAll() ([]*dto.SurveyResponse, error) {
	surveys, err := u.surveyRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.SurveyResponse, len(surveys))
	for i, survey := range surveys {
		responses[i] = convertToResponse(survey)
	}

	return responses, nil
}

func (u *SurveyUsecase) Update(id string, input dto.UpdateSurveyInput, userID string) (*dto.SurveyResponse, error) {
	survey, err := u.surveyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if survey.CreatedBy != userID {
		return nil, apperrors.ErrUnauthorized
	}

	inputVal := reflect.ValueOf(input)
	inputType := inputVal.Type()
	surveyVal := reflect.ValueOf(survey).Elem()

	for i := 0; i < inputVal.NumField(); i++ {
		field := inputVal.Field(i)
		fieldName := inputType.Field(i).Name

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		surveyField := surveyVal.FieldByName(fieldName)
		if !surveyField.IsValid() {
			continue
		}

		if fieldName == "Questions" && !field.IsNil() {
			surveyField.Set(reflect.ValueOf(convertQuestions(input.Questions)))
			continue
		}

		if field.Kind() == reflect.Ptr {
			surveyField.Set(field.Elem())
		}
	}

	if err := u.surveyRepo.Update(survey); err != nil {
		return nil, err
	}

	return convertToResponse(survey), nil
}

func (u *SurveyUsecase) Delete(id string, userID string, role string) error {
	survey, err := u.surveyRepo.GetByID(id)
	if err != nil {
		return err
	}

	if role != string(domain.RoleAdmin) && survey.CreatedBy != userID {
        return apperrors.ErrUnauthorized
    }


	return u.surveyRepo.Delete(id)
}

func (u *SurveyUsecase) GetByUserID(userID string) ([]*dto.SurveyResponse, error) {
	surveys, err := u.surveyRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.SurveyResponse, len(surveys))
	for i, survey := range surveys {
		responses[i] = convertToResponse(survey)
	}

	return responses, nil
}

func convertQuestions(input []dto.QuestionInput) []domain.Question {
	questions := make([]domain.Question, len(input))
	for i, q := range input {
		questions[i] = domain.Question{
			Text:        q.Text,
			Type:        q.Type,
			Required:    q.Required,
			Options:     q.Options,
			Min:         q.Min,
			Max:         q.Max,
			Placeholder: q.Placeholder,
		}
	}
	return questions
}

func convertToResponse(survey *domain.Survey) *dto.SurveyResponse {
	questions := make([]dto.QuestionInput, len(survey.Questions))
	for i, q := range survey.Questions {
		questions[i] = dto.QuestionInput{
			Text:        q.Text,
			Type:        q.Type,
			Required:    q.Required,
			Options:     q.Options,
			Min:         q.Min,
			Max:         q.Max,
			Placeholder: q.Placeholder,
		}
	}

	return &dto.SurveyResponse{
		ID:          survey.ID,
		Title:       survey.Title,
		Description: survey.Description,
		CreatedBy:   survey.CreatedBy,
		Questions:   questions,
		IsActive:    survey.IsActive,
		IsAnonymous: survey.IsAnonymous,
		RequireInfo: survey.RequireInfo,
		StartDate:   survey.StartDate,
		EndDate:     survey.EndDate,
		CreatedAt:   survey.CreatedAt,
		UpdatedAt:   survey.UpdatedAt,
	}
}
