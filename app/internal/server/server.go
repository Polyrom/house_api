package server

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

type Server struct {
	Cfg    *config.Config
	Logger logging.Logger
	Router *mux.Router
	DB     *pgxpool.Pool
}

func (a *Server) ConfigureRouter() {
	ridmw := middleware.NewReqIDMiddleware(a.Logger)
	a.Router.Use(ridmw.DoInMiddle)
	authMwRepo := middleware.NewRepository(a.DB, a.Logger)
	authMwService := middleware.NewService(authMwRepo, a.Logger)
	isAuthMw := middleware.NewAuthMiddleware(authMwService, a.Logger)
	isModerMw := middleware.NewIsModerMiddleware(authMwService, a.Logger)
	urepo := user.NewRepository(a.DB, a.Logger)
	us := user.NewService(urepo, a.Logger)
	ur := user.NewHandler(us, a.Logger)
	ur.Register(a.Router)
	hrepo := house.NewRepository(a.DB, a.Logger)
	hs := house.NewService(hrepo, a.Logger)
	hr := house.NewHandler(isAuthMw, isModerMw, hs, a.Logger)
	hr.Register(a.Router)
	frepo := flat.NewRepository(a.DB, a.Logger)
	fs := flat.NewService(frepo, a.Logger)
	fr := flat.NewHandler(isAuthMw, isModerMw, fs, a.Logger)
	fr.Register(a.Router)
}

func (a *Server) Run() {
	a.Logger.Info("start application")
	srv := &http.Server{
		Addr:         net.JoinHostPort(a.Cfg.Listen.Host, a.Cfg.Listen.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.Logger.Error(err)
		}
	}()
	a.Logger.Infof("server started at :%s", a.Cfg.Listen.Port)

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.Logger.Info("shutting down")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("shutdown server error: %v", err)
	}
	os.Exit(0)
}

func New(cfg *config.Config, logger logging.Logger, router *mux.Router, db *pgxpool.Pool) *Server {
	return &Server{Cfg: cfg, Logger: logger, Router: router, DB: db}
}
