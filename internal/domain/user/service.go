package user

type Service struct {
	repository Repository
}

func NewService(userRepo Repository) *Service {
	return &Service{
		repository: userRepo,
	}
}

func (s *Service) GetUserByID(id string) (*User, error) {
	user, err := s.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ListUsers() ([]*User, error) {
	users, err := s.repository.ListUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
