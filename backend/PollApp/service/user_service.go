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
	req := &store.User{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.NewUser(req.Username, req.Password, req.Email)

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (app *application) GetUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIDStr := ps.ByName("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (app *application) DeleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIDStr := ps.ByName("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	if err := app.store.Users.Delete(r.Context(), userID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
