package main

import (
	"errors"
	"fmt"
	"gaproject.terminator8000.net/internal/data"
	"gaproject.terminator8000.net/internal/validator"
	"net/http"
	"time"
)

func (app *application) expiredToken() {
	fmt.Println("expiredToken started")
	//ticker := time.NewTicker(10 * time.Second)
	for {
		//fmt.Println("")
		users, err := app.models.UserInfo.GetForAllToken()
		if err != nil {
			fmt.Println("GetForAllToken error", err)

			return
		}
		//fmt.Println("users", users)
		for _, user := range users {
			err = app.models.Tokens.DeleteExp(user.ID)
			if err != nil {
				fmt.Println("DeleteExp error")
				return
			}

			token, err := app.models.Tokens.New(user.ID, 5*time.Second, data.ScopeActivation)
			if err != nil {
				fmt.Println("Tokens.New error")
				return
			}

			data := map[string]any{
				"activationToken": token.Plaintext,
				"userID":          user.ID,
			}
			// Send the welcome email, passing in the map above as dynamic data.
			err = app.mailer.Send(user.Email, "new_activation.tmpl", data)
			if err != nil {
				app.logger.PrintError(err, nil)
			}
		}

	}
}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the email and password provided by the client.
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Lookup the user record based on the email address.
	user, err := app.models.UserInfo.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if user is already activated
	if !user.Activated {
		err := app.models.Tokens.Delete(data.ScopeActivation, user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.background(func() {
			data := map[string]any{
				"activationToken": token.Plaintext,
				"userID":          user.ID,
			}
			// Send the welcome email, passing in the map above as dynamic data.
			err = app.mailer.Send(user.Email, "activation.tmpl", data)
			if err != nil {
				app.logger.PrintError(err, nil)
			}
		})

		// Respond with a message indicating account needs activation
		err = app.writeJSON(w, http.StatusCreated, envelope{"message": "Account needs activation. Check your email for the link."}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
