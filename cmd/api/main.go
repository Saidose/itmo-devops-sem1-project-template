package main

import (
	"context"
	"log"
	"net/http"
	"project_sem/internal/config"
	"project_sem/internal/db"
	myhttp "project_sem/internal/http"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	database := db.NewPricesDb(ctx, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost)
	defer database.Close() // закрываем бд
	handlers := myhttp.NewHandlers(database)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/prices", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPrices(w, r)
		case http.MethodPost:
			handlers.PostPrices(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/health", health)
	server := &http.Server{Addr: ":" + cfg.HttpPort, Handler: mux}
	log.Printf("Listening on port: %s", cfg.HttpPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
