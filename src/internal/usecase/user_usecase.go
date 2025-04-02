package usecase

import (
	"time"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/config"
	"survey-project/src/internal/domain"
	"survey-project/src/internal/dto"
	"survey-project/src/pkg/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtConfig        config.JWTConfig
}

func NewUserUsecase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	jwtConfig config.JWTConfig,
) *UserUsecase {
	return &UserUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtConfig:        jwtConfig,
	}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

func (u *UserUsecase) Register(input dto.RegisterInput) error {
	if input.Email == "" || input.Password == "" || input.FullName == "" {
		return apperrors.ErrValidationFailed
	}

	_, err := u.userRepo.GetByEmail(input.Email)
	if err == nil {
		return apperrors.ErrUserAlreadyExists
	}
	if err != apperrors.ErrUserNotFound {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	role := domain.RoleUser
	if input.Role != nil && *input.Role == string(domain.RoleAdmin) {
		role = domain.RoleAdmin
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
		Role:         role,
	}

	return u.userRepo.Create(user)
}

func (u *UserUsecase) Login(input LoginInput) (*TokenPair, error) {
	user, err := u.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUsecase) RefreshToken(token string) (*TokenPair, error) {
	refreshToken, err := u.refreshTokenRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrTokenExpired
	}

	user, err := u.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *UserUsecase) generateAccessToken(user *domain.User) (string, error) {
	claims := &middleware.UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.jwtConfig.AccessDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtConfig.Secret))
}

func (u *UserUsecase) generateRefreshToken(user *domain.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(u.jwtConfig.RefreshDuration).Unix()

	tokenString, err := token.SignedString([]byte(u.jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	if err := u.refreshTokenRepo.UpdateToken(
		user.ID,
		tokenString,
		time.Now().Add(u.jwtConfig.RefreshDuration),
	); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *UserUsecase) GetAll() ([]*dto.UserResponse, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var response []*dto.UserResponse
	for _, user := range users {
		response = append(response, &dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return response, nil
}

func (u *UserUsecase) Delete(id int) error {
	return u.userRepo.Delete(id)
}

func (u *UserUsecase) Update(id int, input *dto.UpdateUserInput, currentUserID int, currentUserRole domain.UserRole) error {
	if currentUserRole != domain.RoleAdmin && currentUserID != id {
		return apperrors.ErrUnauthorized
	}

	updates := make(map[string]interface{})
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.FullName != nil {
		updates["full_name"] = *input.FullName
	}

	return u.userRepo.UpdateFields(id, updates)
}

func (u *UserUsecase) GetByID(id int) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
