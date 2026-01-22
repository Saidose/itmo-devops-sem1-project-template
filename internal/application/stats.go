package application

import "project_sem/internal/domain"

type Stats struct {
	TotalItems      int     `json:"total_items"`
	TotalCategories int     `json:"total_categories"`
	TotalPrice      float64 `json:"total_price"`
}

func GetStats(prices []domain.Price) Stats {
	stats := Stats{0, 0, 0}
	set := make(map[string]bool)
	for _, price := range prices {
		stats.TotalItems++
		set[price.Category] = true
		stats.TotalPrice += price.Price
	}
	stats.TotalCategories = len(set)
	return stats
}
