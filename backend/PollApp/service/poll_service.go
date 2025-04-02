package service

import (
	"PollApp/payload"
	"PollApp/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) CreatePollHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var pollRequest payload.PollRequest
	err := json.NewDecoder(r.Body).Decode(&pollRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pollRequest.CreatorID = 1

	err = app.store.Polls.Create(r.Context(), &pollRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating poll: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) GetPollHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pollIDStr := ps.ByName("pollId")

	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid poll ID: %v", err), http.StatusBadRequest)
		return
	}

	poll, err := app.store.Polls.GetByID(r.Context(), pollID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Poll not found: %v", err), http.StatusNotFound)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, poll); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) ListPollsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	polls, err := app.store.Polls.ListPolls(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching polls: %v", err), http.StatusInternalServerError)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, polls); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) DeletePollHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pollIDStr := ps.ByName("pollId")

	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid poll ID: %v", err), http.StatusBadRequest)
		return
	}

	err = app.store.Polls.Delete(r.Context(), pollID)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, fmt.Sprintf("Poll with ID %d not found", pollID), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error deleting poll: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
