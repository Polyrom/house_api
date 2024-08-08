package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Polyrom/houses_api/internal/config"
	"github.com/Polyrom/houses_api/internal/flat"
	"github.com/Polyrom/houses_api/internal/house"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/internal/user"
	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/gorilla/mux"
)

func main() {
	logger := logging.New()
	cfg := config.Get(logger)
	pg, err := postgres.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("create postgres connection error: %v", err)
	}
	r := mux.NewRouter()
	ridmw := middleware.NewReqIDMiddleware(logger)
	r.Use(ridmw.DoInMiddle)
	authMwRepo := middleware.NewRepository(pg, logger)
	authMwService := middleware.NewService(authMwRepo, logger)
	isAuthMw := middleware.NewAuthMiddleware(authMwService, logger)
	isModerMw := middleware.NewIsModerMiddleware(authMwService, logger)
	urepo := user.NewRepository(pg, logger)
	us := user.NewService(urepo, logger)
	ur := user.NewHandler(us, logger)
	ur.Register(r)
	hrepo := house.NewRepository(pg, logger)
	hs := house.NewService(hrepo, logger)
	hr := house.NewHandler(isAuthMw, isModerMw, hs, logger)
	hr.Register(r)
	frepo := flat.NewRepository(pg, logger)
	fs := flat.NewService(frepo, logger)
	fr := flat.NewHandler(isAuthMw, isModerMw, fs, logger)
	fr.Register(r)
	run(logger, cfg, r)
}

func run(l logging.Logger, c *config.Config, r *mux.Router) {
	l.Info("start application")
	srv := &http.Server{
		Addr:         net.JoinHostPort(c.Listen.Host, c.Listen.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Error(err)
		}
	}()
	l.Infof("server started at :%s", c.Listen.Port)

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l.Info("shutting down")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("shutdown server error: %v", err)
	}
	os.Exit(0)
}
