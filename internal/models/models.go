package models

func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&Product{},
		&Category{},
		&ProductImage{},
		// Add other models here as you create them
	}
}
