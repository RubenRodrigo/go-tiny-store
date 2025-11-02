package auth

import (
	"errors"
	"log"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/token"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/consts"
	"github.com/RubenRodrigo/go-tiny-store/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   user.Repository
	tokenRepo  token.Repository
	jwtManager *jwt.JWTManager
	authConfig consts.AuthConfig
}

func NewService(userRepo user.Repository, tokenRepo token.Repository, jwtManager *jwt.JWTManager) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		tokenRepo:  tokenRepo,
		authConfig: consts.NewAuthConfig(),
	}
}

func (s *Service) RegisterUser(email, username, password, firstName, lastName string) (*user.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if !errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
	}

	if existingUser != nil {
		return nil, apperrors.ErrUserEmailExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Failed to encrypt password. Email: %s, Error: %v",
			email, err)

		return nil, err
	}

	// Create the user
	user := &user.User{
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Authenticate(email, password string) (*user.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrAuthInvalidCredentials
		}

		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, apperrors.ErrAuthInvalidCredentials
	}

	return user, nil
}

func (s *Service) LoginUser(email, password string) (*user.User, string, string, error) {
	// Check if user with email already exists
	user, err := s.Authenticate(email, password)
	if err != nil {
		return nil, "", "", err
	}

	// Generate JWT token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	// Save new refresh token
	_refreshToken := &token.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.authConfig.RefreshTokenTTL),
	}
	if err := s.tokenRepo.Create(_refreshToken); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// LogOutUser implements Service.
func (s *Service) LogOutUser(token string) error {
	if err := s.tokenRepo.DeleteByToken(token); err != nil {
		return err
	}

	return nil
}
