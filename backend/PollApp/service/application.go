package service

import (
	"PollApp/env"
	"PollApp/store"
	"expvar"
	"runtime"

	"go.uber.org/zap"
)

const version = "1.0.0"

var app *application

func GetAppInstance() *application {
	return app
}

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	frontendURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func Start() (string, string) {
	cfg := config{
		addr:        env.GetString("ADDR", ":5020"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:5020"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},

		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Main Database
	// db, err := db.New(
	// 	cfg.db.addr,
	// 	cfg.db.maxOpenConns,
	// 	cfg.db.maxIdleConns,
	// 	cfg.db.maxIdleTime,
	// )
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// defer db.Close()
	// logger.Info("database connection pool established")

	// Rate limiter

	// store := store.NewStorage(db)
	app = &application{
		config: cfg,
		// store:  store,
		logger: logger,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	// expvar.Publish("database", expvar.Func(func() any {
	// 	return db.Stats()
	// }))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	return app.config.addr, app.config.env
}
