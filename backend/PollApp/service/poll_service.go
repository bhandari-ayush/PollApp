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
	pollRequest := &payload.PollRequest{}
	err := json.NewDecoder(r.Body).Decode(&pollRequest)
	if err != nil {
		errResponse := fmt.Errorf("invalid request body: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, pollRequest)

	pollId, err := app.store.Polls.Create(r.Context(), pollRequest)
	if err != nil {
		errResponse := fmt.Errorf("error creating poll: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	app.logger.Info("pollId %d", pollId)

	for _, option := range pollRequest.Options {
		optionId, err := app.store.Polls.CreatePollOption(r.Context(), pollId, option.OptionText)
		app.logger.Infof("optionPollId %d", optionId)
		if err != nil {
			errResponse := fmt.Errorf("error creating poll option: %s", err.Error())
			app.logger.Infof("Error: %s", errResponse.Error())
			app.internalServerError(w, r, errResponse)
			return
		}
	}

	pollRequest.Id = strconv.Itoa(pollId)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pollRequest)
}

func (app *application) GetPollHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pollIdStr := ps.ByName("pollId")
	pollId, err := strconv.Atoi(pollIdStr)
	if err != nil {
		errResponse := fmt.Errorf("invalid poll Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	app.logger.Infof("[%s] pollId:  %s ", r.URL.Path, pollIdStr)
	poll, err := app.store.Polls.GetByID(r.Context(), pollId)
	if err != nil {
		errResponse := fmt.Errorf("poll not found Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	pollOptions, err := app.store.Polls.GetPollOptions(r.Context(), pollId)

	pollResponse := &payload.PollResponse{
		Id:          pollId,
		Description: poll.Description,
		PollOptions: make([]*payload.OptionData, 0),
	}

	for _, option := range pollOptions {
		option := &payload.OptionData{
			OptionId:   option.Id,
			OptionText: option.OptionText,
			VoteCount:  option.VoteCount,
		}
		pollResponse.PollOptions = append(pollResponse.PollOptions, option)
	}

	if err := app.jsonResponse(w, http.StatusOK, pollResponse); err != nil {
		app.logger.Infof("Error: %s", err.Error())
		app.internalServerError(w, r, err)
	}
}

func (app *application) ListPollsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.logger.Infof("[%s] Request recieved ", r.URL.Path)
	polls, err := app.store.Polls.ListPolls(r.Context())
	if err != nil {
		errResponse := fmt.Errorf("error while fetching polls: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, polls); err != nil {
		app.logger.Infof("Error: %s", err.Error())
		app.internalServerError(w, r, err)
	}
}

func (app *application) DeletePollHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pollIdStr := ps.ByName("pollId")
	pollId, err := strconv.Atoi(pollIdStr)
	if err != nil {
		errResponse := fmt.Errorf("invalid poll Id: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.badRequestResponse(w, r, errResponse)
		return
	}

	err = app.store.Votes.DeleteByPollID(r.Context(), pollId)

	if err != nil {
		errResponse := fmt.Errorf("error while deleting poll : %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	err = app.store.Polls.DeletePollOptionById(r.Context(), pollId)
	if err != nil {
		errResponse := fmt.Errorf("error while deleting poll options: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	err = app.store.Polls.Delete(r.Context(), pollId)
	if err != nil {
		if err == store.ErrNotFound {
			errResponse := fmt.Errorf("Poll with ID %d not found : %s", err.Error())
			app.logger.Infof("Error: %s", errResponse.Error())
			app.internalServerError(w, r, errResponse)
			return
		}
		errResponse := fmt.Errorf("error while deleting poll: %s", err.Error())
		app.logger.Infof("Error: %s", errResponse.Error())
		app.internalServerError(w, r, errResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
