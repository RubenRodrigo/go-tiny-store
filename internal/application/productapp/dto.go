package productapp

// ProductFilters represents filters for product queries
type ProductFilters struct {
	CategoryID string
	MinPrice   *float64
	MaxPrice   *float64
	Disabled   *bool
}
