package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
)

type config struct {
	port int
	dev  bool
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config         config
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	templateCache  map[string]*template.Template
}

func main() {
	var cfg config

	// Default flag values for production
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.BoolVar(&cfg.dev, "dev", false, "Development mode")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.Parse()

	// Logger
	h := newSlogHandler(cfg)
	logger := slog.New(h)
	// Create error log for http.Server
	errLog := slog.NewLogLogger(h, slog.LevelError)

	// PostgreSQL
	pool, err := openPool(cfg)
	if err != nil {
		logger.Error("unable to open pgpool", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	// Session manager
	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 12 * time.Hour

	// Template cache
	tc, err := newTemplateCache()
	if err != nil {
		logger.Error("unable to create template cache", slog.Any("error", err))
		os.Exit(1)
	}

	app := &application{
		config:         cfg,
		logger:         logger,
		sessionManager: sm,
		templateCache:  tc,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.port),
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
}

func openPool(cfg config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbpool, err
}

func newSlogHandler(cfg config) slog.Handler {
	if cfg.dev {
		// Development text hanlder
		return tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	}

	// Production use JSON handler with default opts
	return slog.NewJSONHandler(os.Stdout, nil)
}
