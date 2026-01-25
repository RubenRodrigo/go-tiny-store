package models

// AllModels returns all models for schema generation.
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&RefreshToken{},
		&Category{},
		&Product{},
		&ProductImage{},
	}
}
