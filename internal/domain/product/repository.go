package product

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	List() ([]*Product, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// List implements Repository.
func (r *repository) List() ([]*Product, error) {
	panic("unimplemented")
}
