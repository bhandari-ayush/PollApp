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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, pollRequest)

	pollId, err := app.store.Polls.Create(r.Context(), pollRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating poll: %v", err), http.StatusInternalServerError)
		return
	}

	app.logger.Info("pollId %s", pollId)

	for _, option := range pollRequest.Options {
		optionId, err := app.store.Polls.CreatePollOption(r.Context(), pollId, option.OptionText)
		app.logger.Infof("optionPollId %d", optionId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating poll: %v", err), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Invalid poll ID: %v", err), http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] pollId:  %s ", r.URL.Path, pollIdStr)
	poll, err := app.store.Polls.GetByID(r.Context(), pollId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Poll not found: %v", err), http.StatusNotFound)
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
			OptionText: option.OptionText,
			VoteCount:  option.VoteCount,
		}
		pollResponse.PollOptions = append(pollResponse.PollOptions, option)
	}

	if err := app.jsonResponse(w, http.StatusOK, pollResponse); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) ListPollsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.logger.Infof("[%s] Request recieved ", r.URL.Path)
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
	pollIdStr := ps.ByName("pollId")
	pollId, err := strconv.Atoi(pollIdStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid poll ID: %v", err), http.StatusBadRequest)
		return
	}

	app.logger.Infof("[%s] pollID:  %s ", r.URL.Path, pollIdStr)
	err = app.store.Polls.Delete(r.Context(), pollId)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, fmt.Sprintf("Poll with ID %d not found", pollId), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error deleting poll: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
