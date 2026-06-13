package services

import (
	"errors"
	"strings"
	"time"

	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"github.com/rvn14/toys-ecommerce/backend/internal/repositories"
	"github.com/rvn14/toys-ecommerce/backend/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(request dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(request dto.LoginRequest) (*dto.AuthResponse, error)
	GetCurrentUser(userID uint) (*dto.UserResponse, error)
	Logout(tokenID string, expiresAt time.Time) error
}

type authService struct {
	userRepository            repositories.UserRepository
	jwtSecret                 string
	accessTokenExpiresMinutes int
	tokenBlacklist            TokenBlacklist
}

func NewAuthService(userRepository repositories.UserRepository, jwtSecret string, accessTokenExpiresMinutes int, tokenBlacklist TokenBlacklist) AuthService {
	return &authService{
		userRepository:            userRepository,
		jwtSecret:                 jwtSecret,
		accessTokenExpiresMinutes: accessTokenExpiresMinutes,
		tokenBlacklist:            tokenBlacklist,
	}
}

func (s *authService) Register(request dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	existingUser, err := s.userRepository.FindByEmail(email)

	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         request.Name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.RoleCustomer,
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return s.createAuthResponse(user)
}

func (s *authService) Login(request dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(request.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return s.createAuthResponse(user)
}

func (s *authService) GetCurrentUser(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	response := mapUserToResponse(user)
	return &response, nil
}

func (s *authService) Logout(tokenID string, expiresAt time.Time) error {
	if tokenID == "" || expiresAt.IsZero() {
		return ErrUnauthorized
	}
	s.tokenBlacklist.Revoke(tokenID, expiresAt)
	return nil
}

func (s *authService) createAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	token, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		string(user.Role),
		s.jwtSecret,
		s.accessTokenExpiresMinutes,
	)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken: token,
		User:        mapUserToResponse(user),
	}, nil
}

func mapUserToResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}
