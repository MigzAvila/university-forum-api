// Filename: cmd/api/healthcheck.go

package main

import (
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map to display that the server is running
	data := envelope{
		"status": "available",
		"system_info": envelope{
			"environment": app.config.env,
			"version":     version,
		},
	}
	// simulate a delay
	time.Sleep(4 * time.Second)
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		return
	}
}
