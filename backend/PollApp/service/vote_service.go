package service

import (
	"PollApp/payload"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) CreateVoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var voteRequest payload.VoteRequest
	err := json.NewDecoder(r.Body).Decode(&voteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// voteRequest.UserID = 1

	err = app.store.Votes.Create(r.Context(), &voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) UpdateVoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var voteRequest payload.VoteRequest
	err := json.NewDecoder(r.Body).Decode(&voteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// voteRequest.UserID = 1

	err = app.store.Votes.Update(r.Context(), &voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) DeleteVoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var voteRequest payload.VoteRequest
	err := json.NewDecoder(r.Body).Decode(&voteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// voteRequest.UserID = 1

	err = app.store.Votes.Delete(r.Context(), &voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
