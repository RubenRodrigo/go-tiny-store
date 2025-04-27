package service

import (
	"errors"
	"log"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) RegisterUser(email, username, password, firstName, lastName string) (*models.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if !errors.Is(err, apperrors.ErrUserNotFound) {
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
	user := &models.User{
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

func (s *authService) LoginUser(email, password string) (*models.User, string, error) {
	// Check if user with email already exists
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
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
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, "", apperrors.ErrAuthTokenGenerated
	}

	return user, token, nil
}

// generateJWT creates a new JWT token for the authenticated user
func (s *authService) generateJWT(user *models.User) (string, error) {
	// Create claims with user information
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		log.Println("ERROR: Failed to generate token: %w", err)
		return "", err
	}

	return tokenString, nil
}
