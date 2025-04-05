package service

import (
	"PollApp/auth"
	"PollApp/db"
	"PollApp/env"
	"PollApp/store"
	"context"
	"expvar"
	"runtime"
	"time"

	"go.uber.org/zap"
)

const version = "1.0.0"

var app *application

func GetAppInstance() *application {
	return app
}

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	db          dbConfig
	auth        authConfig
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

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

func Start() (string, string) {
	cfg := config{
		addr:        env.GetString("ADDR", ":5020"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:5020"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://root:postgres@10.111.57.52/pollapp?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "1m"),
		},
		auth: authConfig{
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "hello"),
				exp:    time.Hour * 24 * 1,
				iss:    "pollapp",
			},
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()

	defer logger.Sync()

	// Main Database
	db, err := db.ConnectDB(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	// defer db.Close()
	logger.Info("database connection pool established")

	err = store.CreateTables(db, context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	store := store.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)
	app = &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		authenticator: jwtAuthenticator,
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
