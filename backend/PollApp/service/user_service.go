package service

import (
	"PollApp/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	payload := &store.User{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		errResponse := fmt.Errorf("invalid request body: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		app.logger.Infof("Error: %s", err.Error())
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, payload)

	user := store.NewUser(payload.Username, payload.Password, payload.Email)
	id, err := app.store.Users.Create(r.Context(), user)
	if err != nil {
		errResponse := fmt.Errorf("error creating user: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	user.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (app *application) GetUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdStr := ps.ByName("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errResponse := fmt.Errorf("invalid user Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] userId:  %s ", r.URL.Path, userIdStr)

	user, err := app.store.Users.GetByID(r.Context(), userId)
	if err != nil {
		errResponse := fmt.Errorf("error fetching user: %v", err)
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (app *application) DeleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdStr := ps.ByName("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errResponse := fmt.Errorf("invalid user Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] userId:  %s ", r.URL.Path, userIdStr)

	if err := app.store.Users.Delete(r.Context(), userId); err != nil {
		errResponse := fmt.Errorf("error while deleting user: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
