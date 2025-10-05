package auth

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo   user.Repository
	jwtManager *jwt.JWTManager
}

type Service interface {
	RegisterUser(email, username, password, firstName, lastName string) (*user.User, error)
	LoginUser(email, password string) (*user.User, string, error)
	LogOutUser(token string) error
}

func NewService(userRepo user.Repository, jwtManager *jwt.JWTManager) Service {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) RegisterUser(email, username, password, firstName, lastName string) (*user.User, error) {
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

func (s *authService) LoginUser(email, password string) (*user.User, string, error) {
	// Check if user with email already exists
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, "", apperrors.ErrAuthInvalidCredentials
		}

		return nil, "", err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", apperrors.ErrAuthInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, "", apperrors.ErrAuthTokenGenerated
	}

	return user, token, nil
}

// LogOutUser implements AuthService.
func (s *authService) LogOutUser(token string) error {
	return nil
}
