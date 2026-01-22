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
	handlers := myhttp.NewHandlers(database)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v0/prices", handlers.PostPrices)
	mux.HandleFunc("GET /api/v0/prices", handlers.GetPrices)

	server := &http.Server{Addr: ":" + cfg.HttpPort, Handler: mux}
	log.Printf("Listening on port: %s", cfg.HttpPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
