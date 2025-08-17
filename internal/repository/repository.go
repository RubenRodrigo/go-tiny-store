package repository

import "github.com/RubenRodrigo/go-tiny-store/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	ListUsers() ([]*models.User, error)
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type ProductRepository interface {
	List() ([]*models.Product, error)
	GetById(id string) (*models.Product, error)
	Create(product *models.Product) error
	Update(id string, product *models.Product) error
	Delete(id string) error
}

type CategoryRepository interface {
	List() ([]*models.Category, error)
	GetById(id string) (*models.Category, error)
	Create(category *models.Category) error
	Update(id string, category *models.Category) error
	Delete(id string) error
}

type ProductImageRepository interface {
	List() ([]*models.ProductImage, error)
	GetById(id string) (*models.ProductImage, error)
	Create(productImage *models.ProductImage) error
	Update(id string, productImage *models.ProductImage) error
	Delete(id string) error
}
