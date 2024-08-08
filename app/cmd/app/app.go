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
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	logger logging.Logger
	router *mux.Router
	db     *pgxpool.Pool
}

func (a *App) configureRouter() {
	ridmw := middleware.NewReqIDMiddleware(a.logger)
	a.router.Use(ridmw.DoInMiddle)
	authMwRepo := middleware.NewRepository(a.db, a.logger)
	authMwService := middleware.NewService(authMwRepo, a.logger)
	isAuthMw := middleware.NewAuthMiddleware(authMwService, a.logger)
	isModerMw := middleware.NewIsModerMiddleware(authMwService, a.logger)
	urepo := user.NewRepository(a.db, a.logger)
	us := user.NewService(urepo, a.logger)
	ur := user.NewHandler(us, a.logger)
	ur.Register(a.router)
	hrepo := house.NewRepository(a.db, a.logger)
	hs := house.NewService(hrepo, a.logger)
	hr := house.NewHandler(isAuthMw, isModerMw, hs, a.logger)
	hr.Register(a.router)
	frepo := flat.NewRepository(a.db, a.logger)
	fs := flat.NewService(frepo, a.logger)
	fr := flat.NewHandler(isAuthMw, isModerMw, fs, a.logger)
	fr.Register(a.router)
}

func (a *App) run() {
	a.logger.Info("start application")
	srv := &http.Server{
		Addr:         net.JoinHostPort(a.cfg.Listen.Host, a.cfg.Listen.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.logger.Error(err)
		}
	}()
	a.logger.Infof("server started at :%s", a.cfg.Listen.Port)

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.logger.Info("shutting down")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("shutdown server error: %v", err)
	}
	os.Exit(0)
}

func NewApp(cfg *config.Config, logger logging.Logger, router *mux.Router, db *pgxpool.Pool) *App {
	return &App{cfg: cfg, logger: logger, router: router, db: db}
}
