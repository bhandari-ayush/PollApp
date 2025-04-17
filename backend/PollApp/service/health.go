package service

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.logger.Infof("Error: %s", err.Error())
		app.internalServerError(w, r, err)
	}
}
