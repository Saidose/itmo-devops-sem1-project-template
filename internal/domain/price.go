package domain

import "time"

type Price struct {
	ID         int       `csv:"id"`
	Name       string    `csv:"name"`
	Category   string    `csv:"category"`
	Price      float64   `csv:"price"`
	CreateDate time.Time `csv:"create_date"`
}

type PriceCSV struct {
	ID         int     `csv:"id"`
	Name       string  `csv:"name"`
	Category   string  `csv:"category"`
	Price      float64 `csv:"price"`
	CreateDate string  `csv:"create_date"`
}

func PriceConvertDate(prices []Price) []PriceCSV {
	out := make([]PriceCSV, len(prices))
	for i, p := range prices {
		out[i] = PriceCSV{
			ID:         p.ID,
			Name:       p.Name,
			Category:   p.Category,
			Price:      p.Price,
			CreateDate: p.CreateDate.Format("2006-01-02"),
		}
	}
	return out
}
