package service

import (
	"PollApp/store"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=4,max=255"`
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

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)

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
