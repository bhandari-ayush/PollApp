package main

import (
	"PollApp/env"
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

		log.Printf("signal caught %v", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	log.Printf("server has started %s  with env: %s\n", configAddr, environment)

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}
	log.Printf("server has stopped %s  with env: %s\n", configAddr, environment)
	return nil
}

func main() {

	configAddr, environment := service.Start()
	app := service.GetAppInstance()

	apiVersion := env.GetString("API_VERSION", "/v1")

	router := httprouter.New()

	router.GET(apiVersion+"/health", app.HealthCheckHandler)

	// router.GET(apiVersion+"/heatlh/auth", service.WrapHandlerWithParams(app.UpdateVoteHandler))

	// router.POST(apiVersion+"/user/register", app.AuthTokenMiddleware(app.HealthCheckHandler))
	// router.POST(apiVersion + "/user/login")
	// router.POST("/register", app.AuthTokenMiddleware(http.HandlerFunc(app.CreatePollHandler)))

	router.POST(apiVersion+"/poll", app.CreatePollHandler)
	router.GET(apiVersion+"/poll/:pollId", app.GetPollHandler)
	router.GET(apiVersion+"/all/poll/", app.ListPollsHandler)
	router.DELETE(apiVersion+"/poll/:pollId", app.DeletePollHandler)

	router.POST(apiVersion+"/vote", app.CreateVoteHandler)
	router.PUT(apiVersion+"/vote", app.UpdateVoteHandler)

	// router.GET(apiVersion+"/vote/option/:id",app.)
	router.DELETE(apiVersion+"/vote", app.DeleteVoteHandler)

	log.Fatal(run(router, configAddr, environment))
}
