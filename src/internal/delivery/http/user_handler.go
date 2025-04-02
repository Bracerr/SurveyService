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

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	if err := h.userUsecase.Register(input); err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			writeError(w, http.StatusConflict, err)
		case apperrors.ErrValidationFailed:
			writeError(w, http.StatusBadRequest, err)
		case apperrors.ErrInvalidEmail:
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	tokens, err := h.userUsecase.Login(dto.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch err {
		case apperrors.ErrInvalidCredentials:
			writeError(w, http.StatusUnauthorized, err)
		case apperrors.ErrValidationFailed:
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input dto.RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	if input.Token == "" {
		writeError(w, http.StatusBadRequest, apperrors.ErrInvalidToken)
		return
	}

	tokens, err := h.userUsecase.RefreshToken(input.Token)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidToken:
			writeError(w, http.StatusUnauthorized, err)
		case apperrors.ErrTokenExpired:
			writeError(w, http.StatusUnauthorized, err)
		case apperrors.ErrTokenUsed:
			writeError(w, http.StatusUnauthorized, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUsecase.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	user, err := h.userUsecase.GetByID(id)
	if err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	var input dto.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	if err := h.userUsecase.Update(id, &input, claims.UserID, domain.UserRole(claims.Role)); err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		case apperrors.ErrValidationFailed:
			writeError(w, http.StatusBadRequest, err)
		case apperrors.ErrInvalidEmail:
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	if err := h.userUsecase.Delete(id); err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	user, err := h.userUsecase.GetByID(claims.UserID)
	if err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	var input dto.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, apperrors.ErrValidationFailed)
		return
	}

	if err := h.userUsecase.Update(claims.UserID, &input, claims.UserID, domain.UserRole(claims.Role)); err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		case apperrors.ErrValidationFailed:
			writeError(w, http.StatusBadRequest, err)
		case apperrors.ErrInvalidEmail:
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	if err := h.userUsecase.Delete(claims.UserID); err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
