// The entry point for our web app.
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Configuration & settings
const sessionDuration = 24 * time.Hour
const filePath = "./files"
const httpPort = 8080

// The entry point for our server
func main() {
	// Logger init
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Initialize and connect to our database
	initDB()
	defer db.Close()

	// We start with a fresh database every time,
	// so we need to re-create its tables.
	createTables()

	mux := http.NewServeMux()

	// Tell the HTTP server which request should be handled by what function
	setupRoutes(mux)

	// Attach middleware
	httpHandler := panicRecovery(RequestLogging(UserAuth(mux)))

	// Tell the server to start listening
	log.Info("starting the web server at http://localhost" + ":" + strconv.Itoa(httpPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(httpPort), httpHandler))
}

var emptyPageData = NewPageData("", "")

// Define the HTTP routes used by our application
func setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)

		switch request.Method {
		case "GET":
			showPage(response, "index", NewPageData(username, ""))

		default:
			resolveBadRequestMethod(response)
		}
	})

	mux.HandleFunc("/register", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)

		if username != "" {
			showPage(response, "index", NewPageData(username, "Already logged in"))
			return
		}

		switch request.Method {
		case "GET":
			showPage(response, "register", emptyPageData)
		case "POST":
			processRegistration(response, request)

		default:
			resolveBadRequestMethod(response)
		}

	})

	mux.HandleFunc("/login", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)

		if username != "" {
			showPage(response, "index", NewPageData(username, "Already logged in"))
			return
		}

		switch request.Method {
		case "GET":
			showPage(response, "login", emptyPageData)
		case "POST":
			processLoginAttempt(response, request)

		default:
			resolveBadRequestMethod(response)
		}
	})

	mux.HandleFunc("/logout", func(response http.ResponseWriter, request *http.Request) {

		switch request.Method {
		case "GET":
			processLogout(response, request)

		default:
			resolveBadRequestMethod(response)
		}

	})

	mux.HandleFunc("/upload", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)
		data := NewPageData(username, "")

		if username == "" {
			http.Error(response, "Not authorized", http.StatusUnauthorized)
			return
		}

		switch request.Method {
		case "GET":
			showPage(response, "upload", data)
		case "POST":
			processUpload(response, request, username)

		default:
			resolveBadRequestMethod(response)
		}

	})

	mux.HandleFunc("/list", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)

		if username == "" {
			http.Redirect(response, request, "/", http.StatusUnauthorized)
			return
		}

		switch request.Method {
		case "GET":
			listFiles(response, request, username)

		default:
			resolveBadRequestMethod(response)
		}

	})

	mux.HandleFunc("/file/", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)

		if username == "" {
			http.Error(response, "Not authorized", http.StatusUnauthorized)
			return
		}
		switch request.Method {
		case "GET":
			getFile(response, request, username)

		default:
			resolveBadRequestMethod(response)
		}

	})

	mux.HandleFunc("/share", func(response http.ResponseWriter, request *http.Request) {
		username := getUsernameFromCtx(request)
		data := NewPageData(username, "")

		if username == "" {
			http.Error(response, "Not authorized", http.StatusUnauthorized)
			return
		}

		switch request.Method {
		case "GET":
			showPage(response, "share", data)
		case "POST":
			processShare(response, request, username)

		default:
			resolveBadRequestMethod(response)
		}

	})

	// Convenience function for resetting the application's state between tests
	// It should not be used as part of your attacks.
	mux.HandleFunc("/reset", func(response http.ResponseWriter, request *http.Request) {
		resetState()
		fmt.Fprintf(response, "Done")
	})

	// Serve static file such as .css and favicon
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

}

// deals with unrecognised HTTP methods
func resolveBadRequestMethod(response http.ResponseWriter) {
	response.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(response, "Unrecognized method")
}
