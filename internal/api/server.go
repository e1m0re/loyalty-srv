package api

import (
	"context"
	"e1m0re/loyalty-srv/internal/db/migrations"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	appHandler "e1m0re/loyalty-srv/internal/api/handler"
	"e1m0re/loyalty-srv/internal/repository"
	"e1m0re/loyalty-srv/internal/service"
)

type Server struct {
	config     *Config
	httpServer *http.Server
}

func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	db, err := sqlx.Open("pgx", cfg.databaseDSN)
	if err != nil {
		return nil, err
	}

	securityService := service.NewSecurityService(cfg.jwtSecretKey)
	repo := repository.NewRepositories(db)
	services := service.NewServices(repo, securityService)
	handler := appHandler.NewHandler(services)

	srv := &Server{
		config: cfg,
		httpServer: &http.Server{
			Addr:    cfg.serverAddress,
			Handler: handler.NewRouter(),
		},
	}

	return srv, nil
}

func (srv *Server) prepareDB(ctx context.Context) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", srv.config.databaseDSN)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&migrations.Content)
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	return db.Close()
}

func (srv *Server) startHTTPServer(ctx context.Context) error {
	slog.Info(fmt.Sprintf("starting http server at %s", srv.httpServer.Addr))
	err := srv.httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (srv *Server) Start(ctx context.Context) error {
	err := srv.prepareDB(ctx)
	if err != nil {
		return err
	}

	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return srv.startHTTPServer(ctx)
	})

	grp.Go(func() error {
		<-ctx.Done()

		return srv.Shutdown(ctx)
	})

	return grp.Wait()
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
