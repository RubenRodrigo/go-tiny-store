package category

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// List implements CategoryService.
func (s *Service) List() ([]*Category, error) {
	categories, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create implements CategoryService.
func (s *Service) Create(name string) (*Category, error) {
	category := &Category{Name: name}

	if err := s.repository.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Create implements CategoryService.
func (s *Service) Update(name, id string) (*Category, error) {
	category := &Category{Name: name}

	if err := s.repository.Update(id, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (s *Service) Delete(id string) error {

	if err := s.repository.Delete(id); err != nil {
		return err
	}

	return nil
}
