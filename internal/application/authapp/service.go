package authapp

import (
	"fmt"
	"log"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
)

// Service handles authentication-related use cases
type Service struct {
	userRepo               user.Repository
	refreshTokenRepo       user.RefreshTokenRepository
	passwordResetTokenRepo user.PasswordResetTokenRepository
	tokenService           auth.TokenService
	passwordHasher         auth.PasswordHasher
	tokenHasher            auth.TokenHasher
	emailSender            auth.EmailSender
}

// NewService creates a new auth application service
func NewService(
	userRepo user.Repository,
	refreshTokenRepo user.RefreshTokenRepository,
	passwordResetTokenRepo user.PasswordResetTokenRepository,
	tokenService auth.TokenService,
	passwordHasher auth.PasswordHasher,
	tokenHasher auth.TokenHasher,
	emailSender auth.EmailSender,
) *Service {
	return &Service{
		userRepo:               userRepo,
		refreshTokenRepo:       refreshTokenRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		tokenService:           tokenService,
		passwordHasher:         passwordHasher,
		tokenHasher:            tokenHasher,
		emailSender:            emailSender,
	}
}

func (s *Service) SignUp(dto SignUpDTO) (*AuthUserDTO, error) {
	// Hash the password first (fail fast if hashing fails)
	hashedPassword, err := s.passwordHasher.HashPassword(dto.Password)
	if err != nil {
		return nil, apperrors.ErrDatabaseError
	}

	// Create user entity
	newUser := &user.User{
		Email:     dto.Email,
		Username:  dto.Username,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  hashedPassword,
	}

	// Repository handles duplicate email check via unique constraint
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	// Send welcome email (don't fail if email sending fails)
	_ = s.emailSender.SendWelcomeEmail(newUser.Email)

	// Generate tokens
	accessToken, refreshToken, err := s.authenticate(newUser)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &AuthUserDTO{
		ID:           newUser.ID,
		Email:        newUser.Email,
		Username:     newUser.Username,
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) SignIn(dto SignInDTO) (*AuthUserDTO, error) {
	foundUser, err := s.userRepo.GetUserByEmail(dto.Email)
	if err != nil {
		// Convert not found to invalid credentials (don't reveal user existence)
		if err == apperrors.ErrNotFound {
			return nil, apperrors.ErrAuthInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := s.passwordHasher.ComparePassword(foundUser.Password, dto.Password); err != nil {
		return nil, apperrors.ErrAuthInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.authenticate(foundUser)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &AuthUserDTO{
		ID:           foundUser.ID,
		Email:        foundUser.Email,
		Username:     foundUser.Username,
		FirstName:    foundUser.FirstName,
		LastName:     foundUser.LastName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*AuthUserDTO, error) {
	token, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}
	if token.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	if _, err := s.refreshTokenRepo.GetRefreshToken(refreshToken); err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	foundUser, err := s.userRepo.GetUserByID(token.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.DeleteToken(refreshToken); err != nil {
		// Log but don't fail - proceed with new token generation
		fmt.Printf("Warning: failed to delete old refresh token: %v\n", err)
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(foundUser)
	if err != nil {
		return nil, err
	}

	// Save new refresh token
	rt := &user.RefreshToken{
		UserID:    foundUser.ID,
		Token:     newRefreshToken.Token,
		ExpiresAt: newRefreshToken.ExpiresAt,
	}

	if err := s.refreshTokenRepo.SaveToken(rt); err != nil {
		return nil, apperrors.ErrDatabaseError
	}

	// Convert to DTO
	return &AuthUserDTO{
		ID:           foundUser.ID,
		Email:        foundUser.Email,
		Username:     foundUser.Username,
		FirstName:    foundUser.FirstName,
		LastName:     foundUser.LastName,
		AccessToken:  newAccessToken.Token,
		RefreshToken: newRefreshToken.Token,
	}, nil
}

func (s *Service) SignOut(refreshToken string) error {
	return s.refreshTokenRepo.DeleteToken(refreshToken)
}

func (s *Service) ForgotPassword(email string) error {
	foundUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal whether user exists
		return nil
	}

	if err := s.passwordResetTokenRepo.DeleteActiveResetTokens(foundUser.ID); err != nil {
		return nil
	}

	raw, hash, err := s.tokenHasher.GenerateToken()
	if err != nil {
		return nil
	}

	resetToken := &user.PasswordResetToken{
		UserID:    foundUser.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := s.passwordResetTokenRepo.CreateToken(resetToken); err != nil {
		return nil
	}

	resetURL := fmt.Sprintf("https://tiny-store.example.com/reset-password?token=%s", raw)
	log.Printf("Password reset URL for %s: %s", foundUser.Email, resetURL)
	_ = s.emailSender.SendPasswordResetEmail(foundUser.Email, resetURL)
	return nil
}

func (s *Service) ResetPassword(rawToken, newPassword string) error {
	tokenHash := s.tokenHasher.HashToken(rawToken)
	resetToken, err := s.passwordResetTokenRepo.GetTokenByHash(tokenHash)
	if err != nil {
		return apperrors.ErrAuthInvalidResetToken
	}

	hashedPassword, err := s.passwordHasher.HashPassword(newPassword)
	if err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.userRepo.UpdateUserPassword(resetToken.UserID, hashedPassword); err != nil {
		return err
	}

	if err := s.passwordResetTokenRepo.MarkTokenAsUsed(resetToken.TokenHash); err != nil {
		return err
	}

	return s.refreshTokenRepo.DeleteTokensByUserID(resetToken.UserID)
}

// Helper methods

func (s *Service) authenticate(u *user.User) (string, string, error) {
	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(u)
	if err != nil {
		return "", "", err
	}

	// Save refresh token
	rt := &user.RefreshToken{
		UserID:    u.ID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
	}

	if err := s.refreshTokenRepo.SaveToken(rt); err != nil {
		return "", "", apperrors.ErrDatabaseError
	}

	return accessToken.Token, refreshToken.Token, nil
}

func (s *Service) generateTokens(u *user.User) (*auth.GeneratedToken, *auth.GeneratedToken, error) {
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
