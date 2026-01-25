package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/integrations"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
	infraAuth "github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/auth"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo               repository.User
	passwordResetTokenRepo repository.PasswordResetToken
	tokenService           infraAuth.TokenService
	emailService           *integrations.EmailService
}

func NewService(userRepo repository.User, passwordResetTokenRepo repository.PasswordResetToken, tokenService infraAuth.TokenService, emailService *integrations.EmailService) *Service {
	return &Service{
		userRepo:               userRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		tokenService:           tokenService,
		emailService:           emailService,
	}
}

func (s *Service) SignUp(email, username, password, firstName, lastName string) (*AuthUserDTO, error) {
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

	s.emailService.SendWelcomeEmail(newUser.Email)

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

func (s *Service) SignIn(email, password string) (*AuthUserDTO, error) {
	foundUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		// Convert not found to invalid credentials (don't reveal user existence)
		if err == apperrors.ErrNotFound {
			return nil, apperrors.ErrAuthInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return nil, apperrors.ErrAuthInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.authenticate(foundUser)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	authUser := &AuthUserDTO{
		ID:           foundUser.ID,
		Email:        foundUser.Email,
		Username:     foundUser.Username,
		FirstName:    foundUser.FirstName,
		LastName:     foundUser.LastName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authUser, nil
}

func (s *Service) RefreshToken(refreshToken string) (*AuthUserDTO, error) {
	token, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}
	if token.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	if _, err := s.userRepo.GetRefreshToken(refreshToken); err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	foundUser, err := s.userRepo.GetUserById(token.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.DeleteToken(refreshToken); err != nil {
		// Log but don't fail - proceed with new token generation
		// This prevents attack where old tokens remain valid
		fmt.Printf("Warning: failed to delete old refresh token: %v\n", err)
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(foundUser)
	if err != nil {
		return nil, err
	}

	// Save new refresh token
	rt := &models.RefreshToken{
		UserID:    foundUser.ID,
		Token:     newRefreshToken.Token,
		ExpiresAt: newRefreshToken.ExpiresAt,
	}

	if err := s.userRepo.SaveToken(rt); err != nil {
		return nil, apperrors.ErrDatabaseError
	}

	// Convert to DTO
	authUser := &AuthUserDTO{
		ID:           foundUser.ID,
		Email:        foundUser.Email,
		Username:     foundUser.Username,
		FirstName:    foundUser.FirstName,
		LastName:     foundUser.LastName,
		AccessToken:  newAccessToken.Token,
		RefreshToken: newRefreshToken.Token,
	}

	return authUser, nil
}

func (s *Service) SignOut(refreshToken string) error {
	// Simply delete the refresh token
	return s.userRepo.DeleteToken(refreshToken)
}

func (s *Service) ForgotPassword(email string) error {
	foundUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil
	}

	if err := s.passwordResetTokenRepo.DeleteActiveResetTokens(foundUser.ID); err != nil {
		return nil
	}

	raw, hash, err := s.generatePasswordResetToken()
	if err != nil {
		return nil
	}

	resetToken := &models.PasswordResetToken{
		UserID:    foundUser.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := s.passwordResetTokenRepo.CreateToken(resetToken); err != nil {
		return nil
	}

	resetUrl := fmt.Sprintf("https://tiny-store.example.com/reset-password?token=%s", raw)
	log.Printf("Password reset URL for %s: %s", foundUser.Email, resetUrl)
	s.emailService.SendPasswordResetEmail(foundUser.Email, resetUrl)
	return nil
}

func (s *Service) ResetPassword(rawToken, newPassword string) error {
	tokenHash := s.hashToken(rawToken)
	resetToken, err := s.passwordResetTokenRepo.GetTokenByHash(tokenHash)
	if err != nil {
		return apperrors.ErrAuthInvalidResetToken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.userRepo.UpdateUserPassword(resetToken.UserID, string(hashedPassword)); err != nil {
		return err
	}

	if err := s.passwordResetTokenRepo.MarkTokenAsUsed(resetToken.TokenHash); err != nil {
		return err
	}

	return s.userRepo.DeleteTokens(resetToken.UserID)
}

func (s *Service) authenticate(u *models.User) (string, string, error) {
	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(u)
	if err != nil {
		return "", "", err
	}

	// Save refresh token
	rt := &models.RefreshToken{
		UserID:    u.ID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
	}

	if err := s.userRepo.SaveToken(rt); err != nil {
		return "", "", apperrors.ErrDatabaseError
	}

	return accessToken.Token, refreshToken.Token, nil
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

func (s *Service) generatePasswordResetToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(h[:])

	return raw, hash, nil
}

func (s *Service) hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
