package main

import (
	"PollApp/db"
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

	defer db.DisconnectDB()

	router := httprouter.New()

	router.GET(apiVersion+"/health", app.HealthCheckHandler)

	router.POST(apiVersion+"/auth/token", app.CreateTokenHandler)

	router.POST(apiVersion+"/user", app.CreateUserHandler)
	router.GET(apiVersion+"/user/:id", app.GetUserHandler)
	router.DELETE(apiVersion+"/user/:id", app.DeleteUserHandler)

	router.POST(apiVersion+"/poll", app.AuthTokenMiddleware(app.CreatePollHandler))
	router.GET(apiVersion+"/poll/:pollId", app.AuthTokenMiddleware(app.GetPollHandler))
	router.GET(apiVersion+"/all/poll/", app.AuthTokenMiddleware(app.ListPollsHandler))
	router.DELETE(apiVersion+"/poll/:pollId", app.AuthTokenMiddleware(app.DeletePollHandler))
	router.GET(apiVersion+"/option/:optionId/results", app.AuthTokenMiddleware(app.GetOptionVoteUsers))

	router.POST(apiVersion+"/vote", app.AuthTokenMiddleware(app.CreateVoteHandler))
	router.PUT(apiVersion+"/vote", app.AuthTokenMiddleware(app.UpdateVoteHandler))
	router.DELETE(apiVersion+"/vote", app.AuthTokenMiddleware(app.DeleteVoteHandler))

	log.Fatal(run(router, configAddr, environment))
}
