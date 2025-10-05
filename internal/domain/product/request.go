package product

type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=0,max=100"`
}

type UpdateProductRequest struct {
	Name string `json:"name" validate:"required,min=0,max=100"`
}
