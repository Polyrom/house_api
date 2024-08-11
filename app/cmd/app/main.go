package main

import (
	"context"

	"github.com/Polyrom/houses_api/internal/config"
	"github.com/Polyrom/houses_api/internal/server"
	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/gorilla/mux"
)

func main() {
	logger := logging.New()
	cfg := config.Get(logger)
	pg, err := postgres.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		logger.Fatalf("create postgres connection error: %v", err)
	}
	router := mux.NewRouter()
	server := server.New(cfg, logger, router, pg)
	server.ConfigureRouter()
	server.Run()
}
