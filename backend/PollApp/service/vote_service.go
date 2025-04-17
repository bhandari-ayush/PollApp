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
		errResponse := fmt.Errorf("invalid request body: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	voteId, err := app.store.Votes.Create(r.Context(), voteRequest)
	if err != nil {
		errResponse := fmt.Errorf("error creating vote: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
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
		errResponse := fmt.Errorf("invalid request body: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	err = app.store.Votes.Update(r.Context(), voteRequest)
	if err != nil {
		errResponse := fmt.Errorf("error update vote: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) DeleteVoteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	voteRequest := &payload.VoteRequest{}
	err := json.NewDecoder(r.Body).Decode(voteRequest)
	if err != nil {
		errResponse := fmt.Errorf("invalid request body: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, voteRequest)

	err = app.store.Votes.Delete(r.Context(), voteRequest)
	if err != nil {
		errResponse := fmt.Errorf("error while deleting vote : %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) GetOptionVoteUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	optionIdStr := ps.ByName("optionId")
	optionId, err := strconv.Atoi(optionIdStr)
	if err != nil {
		errResponse := fmt.Errorf("invalid option Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] optionID:  %s ", r.URL.Path, optionIdStr)

	users, err := app.store.Votes.GetUsersForOption(r.Context(), optionId)
	if err != nil {
		errResponse := fmt.Errorf("poll not found Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
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
			errResponse := fmt.Errorf("error while fetching user: %s", err.Error())
			app.logger.Infof("Error: %s", errResponse.Error())
			app.internalServerError(w, r, errResponse)
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
