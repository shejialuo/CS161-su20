package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type key int

const userKey key = 0

// RequestLogging is a HTTP middleware that logs each incoming http request
func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// pass the request to the next handler
		next.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"path":           r.RequestURI,
			"execution_time": time.Since(startTime).String(),
			"remote_addr":    r.RemoteAddr,
		}).Info("received a new http request")
	})
}

// UserAuth Middleware
func UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		// get the session token cookie
		cookie, err := request.Cookie("session_token")
		_, _ = cookie, err

		if err != nil {
			next.ServeHTTP(w, request)
			return
		}

		// Extract the session token value from the cookie
		sessionToken := cookie.Value
		// empty assignment to suppress unused variable warning
		_ = sessionToken

		//////////////////////////////////
		// BEGIN TASK 1: YOUR CODE HERE
		//////////////////////////////////

		// TODO: look up the session token in the database

		// TODO: make sure the session token exists (i.e. your query returned something)

		// TODO: assign the results of your query to some variables

		// TODO: check that the session token has not expired
		// hint: time.Unix, time.Now, and x.Before(y) may be useful here

		// TODO: if the session token is valid, run the following line of code,
		//       with username assigned to the appropriate value:
		// request = request.WithContext(context.WithValue(request.Context(), userKey, username))

		// TODO: if the session token is invalid, run the following line of code:
		// next.ServeHTTP(w, request)

		//////////////////////////////////
		// END TASK 1: YOUR CODE HERE
		//////////////////////////////////
	})
}

// This method extracts the username from the context of the HTTP request.
// It returns "" when the value is not present in the context.
func getUsernameFromCtx(request *http.Request) string {
	var username string
	usernameCtxValue := request.Context().Value(userKey)
	if usernameCtxValue == nil {
		username = ""
	} else {
		username = usernameCtxValue.(string)
	}
	return username
}

// HTTP panic recovery
// When the handler panics, the program will recover from the panic and logs the error to the console.
// This avoids the termination of the server program.
func panicRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				// To avoid 'superfluous response.WriteHeader call' error
				if rw.Header().Get("Content-Type") == "" {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		handler.ServeHTTP(rw, rq)
	})
}
