package category_handler

type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=0,max=100"`
}
