package service

import (
	"PollApp/store"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type CreateUserTokenPayload struct {
	Id    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email,max=255"`
}

func (app *application) CreateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	payload := &CreateUserTokenPayload{}
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Infof("[%s] Payload %+v ", r.URL.Path, payload)

	userID, err := strconv.Atoi(payload.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)

	app.logger.Infof("user %+v ", user)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if strings.Compare(user.Email, payload.Email) != 0 {
		app.unauthorizedErrorResponse(w, r, store.ErrNotMatch)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	app.logger.Infof("[%s] Payload token %+v ", r.URL.Path, payload)
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, map[string]string{"token": token}); err != nil {
		app.internalServerError(w, r, err)
	}
}
