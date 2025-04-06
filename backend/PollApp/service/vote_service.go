package service

import (
	"PollApp/payload"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) CreateVoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	voteRequest := &payload.VoteRequest{}
	err := json.NewDecoder(r.Body).Decode(voteRequest)
	if err != nil {
		app.logger.Infof("error %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	voteId, err := app.store.Votes.Create(r.Context(), voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating vote: %v", err), http.StatusInternalServerError)
		return
	}

	voteRequest.Id = voteId
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(voteRequest)
}

func (app *application) UpdateVoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	voteRequest := &payload.VoteRequest{}
	err := json.NewDecoder(r.Body).Decode(voteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	err = app.store.Votes.Update(r.Context(), voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) DeleteVoteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	voteRequest := &payload.VoteRequest{}
	err := json.NewDecoder(r.Body).Decode(voteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	err = app.store.Votes.Delete(r.Context(), voteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) GetOptionVoteUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	optionIdStr := ps.ByName("optionId")
	optionId, err := strconv.Atoi(optionIdStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid option ID: %v", err), http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] optionID:  %s ", r.URL.Path, optionIdStr)

	users, err := app.store.Votes.GetUsersForOption(r.Context(), optionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Poll not found: %v", err), http.StatusNotFound)
		return
	}

	optionUserResponse := &payload.OptionUserResponse{
		OptionId:  optionIdStr,
		VoteCount: len(users),
		Users:     make([]*payload.User, 0),
	}

	for _, userId := range users {
		userData, err := app.store.Users.GetByID(r.Context(), userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
			return
		}
		user := &payload.User{
			Name:  userData.Username,
			Email: userData.Email,
		}
		optionUserResponse.Users = append(optionUserResponse.Users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(optionUserResponse)

}
