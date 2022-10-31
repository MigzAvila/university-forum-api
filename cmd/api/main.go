// File: cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"universityforum.miguelavila.net/internals/data"
	"universityforum.miguelavila.net/internals/jsonlog"
)

// App Version
const version = "1.0.0"

// App config
type config struct {
	port int
	env  string // dev, stg, prd, etc...
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps    float64 //request per second
		burst  int
		enable bool
	}
}

// dependencies injections
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config
	//read in the flag that are needed to populate the config ~ flag for using as extra cmd
	flag.IntVar(&cfg.port, "port", 4000, "API port")
	flag.StringVar(&cfg.env, "env", "dev", "(dev | stg | prd)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("FORUM_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-open-conns", 25, "PostgreSQL max idle open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-open-time", "15m", "PostgreSQL max connections idle time")

	// Flag for rate limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enable, "limiter-enable", true, "Enable Rate Limiter")

	flag.Parse()

	//create a logger ~ use := for undeclared var
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//create the connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// memory leak prevent
	defer db.Close()

	// log successful connection
	logger.PrintInfo("database connection pool established", nil)

	//create instances of out api
	app := &application{
		config: cfg,
		logger: logger,
		models: *data.NewModels(db),
	}

	// Call app.serve() to start the server
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

// openDB return a *sql.DB instance
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// create a context with a 5 section timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
