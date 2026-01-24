package auth

import (
	"errors"
	"log"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/consts"
	"github.com/RubenRodrigo/go-tiny-store/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService UserServicePort
	jwtManager  *jwt.JWTManager
	authConfig  consts.AuthConfig
}

func NewService(userService UserServicePort, jwtManager *jwt.JWTManager) Service {
	return &service{
		userService: userService,
		jwtManager:  jwtManager,
		authConfig:  consts.NewAuthConfig(),
	}
}

func (s *service) RegisterUser(email, username, password, firstName, lastName string) (*AuthUserDTO, error) {
	// Check if user with email already exists
	existingUser, err := s.userService.GetByEmail(email)
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

	user, err := s.userService.Create(RegisterUserDto{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) authenticate(email, password string) (*AuthUserDTO, error) {
	user, err := s.userService.GetByEmail(email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrAuthInvalidCredentials
		}

		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, apperrors.ErrAuthInvalidCredentials
	}

	return user, nil
}

func (s *service) LoginUser(email, password string) (*AuthUserDTO, string, string, error) {
	// Check if authUser with email already exists
	authUser, err := s.authenticate(email, password)
	if err != nil {
		return nil, "", "", err
	}

	// Generate JWT token
	accessToken, err := s.jwtManager.GenerateAccessToken(authUser.ID, authUser.Email, authUser.Username)
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	expiresAt := time.Now().Add(s.authConfig.RefreshTokenTTL)
	if err := s.userService.IssueRefreshToken(expiresAt, authUser.ID, refreshToken); err != nil {
		return nil, "", "", err
	}

	return authUser, accessToken, refreshToken, nil
}

func (s *service) RefreshToken(token string) (*AuthUserDTO, string, string, error) {
	// Check if authUser with email already exists
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, "", "", err
	}

	userId, err := jwt.ExtractStringClaim(claims, "user_id")
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenInvalid
	}

	authUser, err := s.userService.GetById(userId)
	if err != nil {
		return nil, "", "", apperrors.ErrNotFound
	}

	// TODO: Delete Previous refresh_token

	// Generate JWT token
	accessToken, err := s.jwtManager.GenerateAccessToken(authUser.ID, authUser.Email, authUser.Username)
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenGenerated
	}

	expiresAt := time.Now().Add(s.authConfig.RefreshTokenTTL)
	if err := s.userService.IssueRefreshToken(expiresAt, authUser.ID, refreshToken); err != nil {
		return nil, "", "", err
	}

	return authUser, accessToken, refreshToken, nil
}

// LogOutUser implements Service.
func (s *service) LogOutUser(token string) error {
	if err := s.userService.RevokeToken(token); err != nil {
		return err
	}

	return nil
}
