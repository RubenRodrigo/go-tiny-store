package authadapter

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
)

var _ auth.UserServicePort = (*Adapter)(nil)

type Adapter struct {
	svc user.Service
}

func NewAdapter(svc user.Service) *Adapter {
	return &Adapter{svc}
}

func (a *Adapter) Create(input auth.RegisterUserDto) (*auth.AuthUserDTO, error) {
	u, err := a.svc.Create(user.CreateUserInput{
		Email:        input.Email,
		Username:     input.Username,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: input.Password,
	})

	if err != nil {
		return nil, err
	}

	return &auth.AuthUserDTO{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PasswordHash: u.Password,
	}, nil
}

func (a *Adapter) GetByEmail(email string) (*auth.AuthUserDTO, error) {
	u, err := a.svc.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return &auth.AuthUserDTO{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PasswordHash: u.Password,
	}, nil
}

func (a *Adapter) IssueRefreshToken(time time.Time, userId, tokenString string) error {
	_, err := a.svc.IssueRefreshToken(time, userId, tokenString)

	return err
}

func (a *Adapter) RevokeToken(token string) error {
	return a.svc.RevokeToken(token)
}
