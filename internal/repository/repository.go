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
	Update(product *models.Product) error
	Delete(product *models.Product) error
}

type CategoryRepository interface {
	List() ([]*models.Category, error)
	GetById(id string) (*models.Category, error)
	Save(category *models.Category) error
	Delete(category *models.Category) error
}

type ProductImageRepository interface {
	List() ([]*models.ProductImage, error)
	GetById(id string) (*models.ProductImage, error)
	Create(productImage *models.ProductImage) error
	Update(productImage *models.ProductImage) error
	Delete(productImage *models.ProductImage) error
}
