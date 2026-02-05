package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/recchia/greenlight/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn                string
		maxOpenConnections int
		maxIdleConnections int
		maxIdleTime        time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	flag.IntVar(&cfg.port, "port", 4000, "port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "Database connection string")
	flag.IntVar(&cfg.db.maxOpenConnections, "max-open-connections", 25, "Maximum number of open connections to the database")
	flag.IntVar(&cfg.db.maxIdleConnections, "max-idle-connections", 25, "Maximum number of idle connections in the pool")
	flag.DurationVar(&cfg.db.maxIdleTime, "max-idle-time", 15*time.Minute, "Maximum duration that a connection can be idle before being closed")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Requests per second limit for the rate limiter")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Burst size for the rate limiter")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiting")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConnections)
	db.SetMaxIdleConns(cfg.db.maxIdleConnections)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
