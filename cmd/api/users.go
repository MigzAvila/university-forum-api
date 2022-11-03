// Filename: cmd/api/forum.go

package main

import (
	"errors"
	"net/http"

	"universityforum.miguelavila.net/internals/data"
	"universityforum.miguelavila.net/internals/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// hold data from the request body
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body into the input struct
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badResquestReponse(w, r, err)
		return
	}

	// Copy data to a new struct
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Generate the password hash from the password the user provided
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Perform validation

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// insert the data in the database
	err = app.models.User.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.background(func() {
		// Send email to new user
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	// write a 202 status code indicating that the user has been Accepted but not created successfully
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
