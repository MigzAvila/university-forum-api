// Filename: cmd/api/forum.go

package main

import (
	"net/http"
)

func (app *application) signupUserHandler(w http.ResponseWriter, r *http.Request) {

	app.logger.Println("Creating user...")
}
