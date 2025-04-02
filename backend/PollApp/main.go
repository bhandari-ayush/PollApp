package main

import (
	"PollApp/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func run(r http.Handler, configAddr string, environment string) error {
	srv := &http.Server{
		Addr:         configAddr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	log.Printf("server has started: %s  with env: %s\n", configAddr, environment)

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	// log.Printf("server has started: %s  with env: %s\n", configAddr, environment)

	return nil
}

func main() {

	configAddr, environment := service.Start()

	app := service.GetAppInstance()

	log.Printf("application %+v", app)

	router := httprouter.New()

	router.GET("/v1/health", app.HealthCheckHandler)

	log.Fatal(run(router, configAddr, environment))
}
