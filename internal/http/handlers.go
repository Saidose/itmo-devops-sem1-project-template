package http

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"project_sem/internal/db"
	"project_sem/internal/domain"

	"github.com/gocarina/gocsv"
)

type Handlers struct {
	db *db.PricesDB
}

func NewHandlers(db *db.PricesDB) *Handlers {
	return &Handlers{db}
}

func (h *Handlers) PostPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var prices []domain.Price

	for _, f := range zr.File {
		if f.Name == "data.csv" {
			rc, err := f.Open()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer rc.Close()

			if err := gocsv.Unmarshal(rc, &prices); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	priceStats, err := h.db.InsertPrices(ctx, prices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priceStats)
}

func (h *Handlers) GetPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prices, err := h.db.GetAll(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "stats/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=response.zip")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	f, err := zipWriter.Create("data.csv")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := gocsv.Marshal(prices, f); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
