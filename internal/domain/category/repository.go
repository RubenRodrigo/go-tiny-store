package category

// Repository defines the interface for category persistence operations
type Repository interface {
	ListCategories() ([]*Category, error)
	GetCategoryByID(id string) (*Category, error)
	CreateCategory(category *Category) error
	UpdateCategory(id string, category *Category) error
	DeleteCategory(id string) error
}
