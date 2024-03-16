package main

import (
	"context"
	"errors"
	"github.com/ast3am/VKintern-movies/api/handlers"
	"github.com/ast3am/VKintern-movies/internal/config"
	"github.com/ast3am/VKintern-movies/internal/db"
	"github.com/ast3am/VKintern-movies/internal/service"
	"github.com/ast3am/VKintern-movies/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//@title VKintern api doc
//@version 1.0

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	cfg := config.GetConfig("config/config.yml")
	log := logging.GetLogger(cfg.LogLevel, os.Stdout)
	db, err := db.NewClient(ctx, cfg, log)
	if err != nil {
		log.FatalMsg("", err)
	}
	defer db.Close(ctx)

	mux := http.NewServeMux()
	service := service.NewService(db, log)
	handler := handlers.NewHandler(service, log)
	handler.RegisterHandlers(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.FatalMsg("listen: %s\n", err)
		}
	}()

	log.InfoMsg("service is running")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
