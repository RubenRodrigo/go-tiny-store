package product

// Filters represents filtering criteria for product queries (domain value object)
type Filters struct {
	CategoryID string
	MinPrice   *float64
	MaxPrice   *float64
	Disabled   *bool
}
