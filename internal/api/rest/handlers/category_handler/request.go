package category_handler

type SaveCategoryRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required,min=0,max=100"`
}
