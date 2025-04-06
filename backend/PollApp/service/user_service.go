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
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, payload)

	user := store.NewUser(payload.Username, payload.Password, payload.Email)
	id, err := app.store.Users.Create(r.Context(), user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Invalid user Id: %v", err), http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] userId:  %s ", r.URL.Path, userIdStr)

	user, err := app.store.Users.GetByID(r.Context(), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (app *application) DeleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdStr := ps.ByName("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] userId:  %s ", r.URL.Path, userIdStr)

	if err := app.store.Users.Delete(r.Context(), userId); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
