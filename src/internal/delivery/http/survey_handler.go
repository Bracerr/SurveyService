package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/domain"
	"survey-project/src/internal/dto"
	"survey-project/src/internal/usecase"
	"survey-project/src/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

type SurveyHandler struct {
	surveyUsecase *usecase.SurveyUsecase
}

func NewSurveyHandler(surveyUsecase *usecase.SurveyUsecase) *SurveyHandler {
	return &SurveyHandler{
		surveyUsecase: surveyUsecase,
	}
}

func (h *SurveyHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	var input dto.CreateSurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	survey, err := h.surveyUsecase.Create(input, strconv.Itoa(claims.UserID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, survey)
}

func (h *SurveyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	id := chi.URLParam(r, "id")
	survey, err := h.surveyUsecase.GetByID(id, strconv.Itoa(claims.UserID), claims.Role)
	if err != nil {
		switch err {
		case apperrors.ErrSurveyNotFound:
			writeError(w, http.StatusNotFound, err)
		case apperrors.ErrUnauthorized:
			writeError(w, http.StatusForbidden, err)
		default:
			writeError(w, http.StatusInternalServerError, apperrors.ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusOK, survey)
}

func (h *SurveyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	surveys, err := h.surveyUsecase.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, apperrors.ErrInternalServer)
		return
	}

	writeJSON(w, http.StatusOK, surveys)
}

func (h *SurveyHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	id := chi.URLParam(r, "id")
	var input dto.UpdateSurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	survey, err := h.surveyUsecase.Update(id, input, strconv.Itoa(claims.UserID))
	if err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		case apperrors.ErrUnauthorized:
			writeError(w, http.StatusForbidden, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, survey)
}

func (h *SurveyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.surveyUsecase.Delete(id, strconv.Itoa(claims.UserID), claims.Role); err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		case apperrors.ErrUnauthorized:
			writeError(w, http.StatusForbidden, err)
		case apperrors.ErrSurveyNotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SurveyHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	if claims.Role != string(domain.RoleAdmin) && claims.UserID != userID {
		writeError(w, http.StatusForbidden, apperrors.ErrUnauthorized)
		return
	}

	surveys, err := h.surveyUsecase.GetByUserID(strconv.Itoa(userID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, apperrors.ErrInternalServer)
		return
	}

	writeJSON(w, http.StatusOK, surveys)
}

func (h *SurveyHandler) GetMy(w http.ResponseWriter, r *http.Request) {
    claims, err := middleware.GetUserFromContext(r.Context())
    if err != nil {
        writeError(w, http.StatusUnauthorized, err)
        return
    }

    surveys, err := h.surveyUsecase.GetByUserID(strconv.Itoa(claims.UserID))
    if err != nil {
        writeError(w, http.StatusInternalServerError, apperrors.ErrInternalServer)
        return
    }

    writeJSON(w, http.StatusOK, surveys)
}
