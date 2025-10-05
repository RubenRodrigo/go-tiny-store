package category

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=0,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=0,max=100"`
}
