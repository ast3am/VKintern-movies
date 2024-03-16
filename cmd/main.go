package main

import (
	"context"
	"github.com/ast3am/VKintern-movies/api/handlers"
	"github.com/ast3am/VKintern-movies/internal/config"
	"github.com/ast3am/VKintern-movies/internal/db"
	"github.com/ast3am/VKintern-movies/internal/service"
	"github.com/ast3am/VKintern-movies/pkg/logging"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()
	cfg := config.GetConfig("config/config.yml")
	log := logging.GetLogger(cfg.LogLevel, os.Stdout)
	db, err := db.NewClient(ctx, cfg, log)
	if err != nil {
		log.FatalMsg("", err)
	}
	defer db.Close(ctx)

	service := service.NewService(db, log)
	handler := handlers.NewHandler(service, log)
	handler.RegisterHandlers(mux)
	log.Print("Starting server on :8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.FatalMsg("Failed to start server: %v", err)
	}
}
