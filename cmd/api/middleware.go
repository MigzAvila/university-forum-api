// Filename: cmd/api/middleware.go

package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})

}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// limit := rate.NewLimiter(2, 4)
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	if !limit.Allow() {
	// 		app.rateLimitExceededResponse(w, r)
	// 		return
	// 	}
	// 	next.ServeHTTP(w, r)
	// })

	// creaet a client type
	type Client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)
	// launch a background goroutine that remove old clients
	// from the clients map every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock before starting to clean up
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enable {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			// Lock
			mu.Lock()
			// check if the IP is in the map
			if _, found := clients[ip]; !found {
				clients[ip] = &Client{limiter: rate.NewLimiter(
					rate.Limit(app.config.limiter.rps),
					app.config.limiter.burst,
				)}
			}

			// update lastSeen
			clients[ip].lastSeen = time.Now()

			// check if the request is allowed
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		} // end of enable conditional
		next.ServeHTTP(w, r)
	})

}
