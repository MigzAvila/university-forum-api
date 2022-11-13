// Filename: cmd/api/forum.go

package main

import (
	"net/http"
)

func (app *application) getForumHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.PrintInfo("dummy function...", nil)
}
