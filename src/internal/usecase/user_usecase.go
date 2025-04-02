package usecase

import (
	"regexp"
	"time"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/config"
	"survey-project/src/internal/domain"
	"survey-project/src/internal/dto"
	"survey-project/src/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type UserUsecase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtConfig        config.JWTConfig
}

func NewUserUsecase(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, jwtConfig config.JWTConfig) *UserUsecase {
	return &UserUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtConfig:        jwtConfig,
	}
}

func (u *UserUsecase) Register(input dto.RegisterInput) error {
	if input.Email == "" || input.Password == "" || input.FullName == "" {
		return apperrors.ErrValidationFailed
	}

	if !isEmailValid(input.Email) {
		return apperrors.ErrInvalidEmail
	}

	if _, err := u.userRepo.GetByEmail(input.Email); err == nil {
		return apperrors.ErrUserAlreadyExists
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

func (u *UserUsecase) Login(input dto.LoginInput) (*dto.TokenPair, error) {
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

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUsecase) RefreshToken(token string) (*dto.TokenPair, error) {
	refreshToken, err := u.refreshTokenRepo.GetByToken(token)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
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

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *UserUsecase) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"exp":     time.Now().Add(u.jwtConfig.AccessDuration).Unix(),
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

	err = u.refreshTokenRepo.UpdateToken(user.ID, tokenString, time.Now().Add(u.jwtConfig.RefreshDuration))
	if err != nil {
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

	if input.Email != nil && !isEmailValid(*input.Email) {
		return apperrors.ErrInvalidEmail
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
