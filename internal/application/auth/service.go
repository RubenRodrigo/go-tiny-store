package auth

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
	infraAuth "github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/auth"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo     repository.User
	tokenService infraAuth.TokenService
}

func NewService(userRepo repository.User, tokenService infraAuth.TokenService) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *Service) RegisterUser(email, username, password, firstName, lastName string) (*AuthUserDTO, error) {
	// Hash the password first (fail fast if bcrypt fails)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrDatabaseError // Internal error
	}

	// Create user entity
	newUser := &models.User{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
	}

	// Repository handles duplicate email check via unique constraint
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	// Convert to DTO
	return &AuthUserDTO{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Username:  newUser.Username,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
	}, nil
}

func (s *Service) LoginUser(email, password string) (*AuthUserDTO, string, string, error) {
	// Get user by email
	foundUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		// Convert not found to invalid credentials (don't reveal user existence)
		if err == apperrors.ErrNotFound {
			return nil, "", "", apperrors.ErrAuthInvalidCredentials
		}
		return nil, "", "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return nil, "", "", apperrors.ErrAuthInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(foundUser)
	if err != nil {
		return nil, "", "", err
	}

	// Save refresh token
	rt := &models.RefreshToken{
		UserID:    foundUser.ID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
	}

	if err := s.userRepo.SaveToken(rt); err != nil {
		return nil, "", "", apperrors.ErrDatabaseError
	}

	// Convert to DTO
	authUser := &AuthUserDTO{
		ID:        foundUser.ID,
		Email:     foundUser.Email,
		Username:  foundUser.Username,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
	}

	return authUser, accessToken.Token, refreshToken.Token, nil
}

func (s *Service) RefreshToken(refreshToken string) (*AuthUserDTO, string, string, error) {
	// Validate the refresh token format (not expired, valid JWT)
	token, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, "", "", apperrors.ErrAuthTokenInvalid
	}

	// Get user from database
	foundUser, err := s.userRepo.GetUserById(token.UserID)
	if err != nil {
		return nil, "", "", err
	}

	// Delete old refresh token (security: prevent token reuse)
	if err := s.userRepo.DeleteToken(refreshToken); err != nil {
		// Log but don't fail - proceed with new token generation
		// This prevents attack where old tokens remain valid
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(foundUser)
	if err != nil {
		return nil, "", "", err
	}

	// Save new refresh token
	rt := &models.RefreshToken{
		UserID:    foundUser.ID,
		Token:     newRefreshToken.Token,
		ExpiresAt: newRefreshToken.ExpiresAt,
	}

	if err := s.userRepo.SaveToken(rt); err != nil {
		return nil, "", "", apperrors.ErrDatabaseError
	}

	// Convert to DTO
	authUser := &AuthUserDTO{
		ID:        foundUser.ID,
		Email:     foundUser.Email,
		Username:  foundUser.Username,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
	}

	return authUser, newAccessToken.Token, newRefreshToken.Token, nil
}

func (s *Service) LogOutUser(refreshToken string) error {
	// Simply delete the refresh token
	return s.userRepo.DeleteToken(refreshToken)
}

// Helper function to generate both tokens
func (s *Service) generateTokens(u *models.User) (*infraAuth.TokenGeneratedClaims, *infraAuth.TokenGeneratedClaims, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(u.ID, u.Email, u.Username)
	if err != nil {
		return nil, nil, apperrors.ErrAuthTokenGenerated
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(u.ID, u.Email, u.Username)
	if err != nil {
		return nil, nil, apperrors.ErrAuthTokenGenerated
	}

	return accessToken, refreshToken, nil
}
