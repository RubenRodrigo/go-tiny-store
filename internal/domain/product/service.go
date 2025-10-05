package product

type Service struct {
	productRepo Repository
}

func NewService(productRepo Repository) *Service {
	return &Service{
		productRepo: productRepo,
	}
}

// Create implements ProductService.
func (p *Service) Create(name string) (*Product, error) {
	panic("unimplemented")
}

// Delete implements ProductService.
func (p *Service) Delete(id string) error {
	panic("unimplemented")
}

// List implements ProductService.
func (p *Service) List() ([]*Product, error) {
	panic("unimplemented")
}

// Update implements ProductService.
func (p *Service) Update(name string, id string) (*Product, error) {
	panic("unimplemented")
}
