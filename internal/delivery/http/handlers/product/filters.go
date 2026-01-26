package product

import (
	"net/http"
	"strconv"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/productapp"
)

// ParseFilters parses product filters from query parameters
func ParseFilters(r *http.Request) productapp.ProductFilters {
	filters := productapp.ProductFilters{}

	// Category filter
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		filters.CategoryID = categoryID
	}

	// Price range filters
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	// Disabled filter
	if disabledStr := r.URL.Query().Get("disabled"); disabledStr != "" {
		if disabled, err := strconv.ParseBool(disabledStr); err == nil {
			filters.Disabled = &disabled
		}
	}

	return filters
}
