// Filename: cmd/api/forum.go

package main

import (
	"net/http"
)

func (app *application) dummy(w http.ResponseWriter, r *http.Request) {

	app.logger.Println("dummy function...")
}
