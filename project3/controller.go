// Main logic/functionality for the web application.
// This is where you need to implement your own server.
package main

// Reminder that you're not allowed to import anything that isn't part of the Go standard library.
// This includes golang.org/x/
import (
	"database/sql"
	"fmt"
	_ "io/ioutil"
	"net/http"
	_ "os"
	_ "path/filepath"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

func processRegistration(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	// Check if username already exists
	row := db.QueryRow("SELECT username FROM users WHERE username = ?", username)
	var savedUsername string
	err := row.Scan(&savedUsername)
	if err != sql.ErrNoRows {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, "username %s already exists", savedUsername)
		return
	}

	// Generate salt
	const saltSizeBytes = 16
	salt, err := randomByteString(saltSizeBytes)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, err.Error())
		return
	}

	hashedPassword := hashPassword(password, salt)

	_, err = db.Exec("INSERT INTO users VALUES (NULL, ?, ?, ?)", username, hashedPassword, salt)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, err.Error())
		return
	}

	// Set a new session cookie
	initSession(response, username)

	// Redirect to next page
	http.Redirect(response, request, "/", http.StatusFound)
}

func processLoginAttempt(response http.ResponseWriter, request *http.Request) {
	// Retrieve submitted values
	username := request.FormValue("username")
	password := request.FormValue("password")

	row := db.QueryRow("SELECT password, salt FROM users WHERE username = ?", username)

	// Parse database response: check for no response or get values
	var encodedHash, encodedSalt string
	err := row.Scan(&encodedHash, &encodedSalt)
	if err == sql.ErrNoRows {
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(response, "unknown user")
		return
	} else if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, err.Error())
		return
	}

	// Hash submitted password with salt to allow for comparison
	submittedPassword := hashPassword(password, encodedSalt)

	// Verify password
	if submittedPassword != encodedHash {
		fmt.Fprintf(response, "incorrect password")
		return
	}

	// Set a new session cookie
	initSession(response, username)

	// Redirect to next page
	http.Redirect(response, request, "/", http.StatusFound)
}

func processLogout(response http.ResponseWriter, request *http.Request) {
	// get the session token cookie
	cookie, err := request.Cookie("session_token")
	// empty assignment to suppress unused variable warning
	_, _ = cookie, err

	//////////////////////////////////
	// BEGIN TASK 2: YOUR CODE HERE
	//////////////////////////////////

	// optional: use this line to get the username of the current user
	// username := getUsernameFromCtx(request)

	// TODO: clear the session token cookie in the user's browser
	// hint: to clear a cookie, set its MaxAge to -1

	// TODO: delete the session from the database

	//////////////////////////////////
	// END TASK 2: YOUR CODE HERE
	//////////////////////////////////

	// redirect to the homepage
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

func processUpload(response http.ResponseWriter, request *http.Request, username string) {

	// BEGIN TASK 3: YOUR CODE HERE
	// you are recommended to use this (const filePath = "./files") as your base directory for file storage
	// you can create sub directories too

	// replace this statement
	fmt.Fprintf(response, "placeholder")

	// END TASK 3: YOUR CODE HERE
}

// fileInfo helps you pass information to the template
type fileInfo struct {
	Filename  string
	FileOwner string
	FilePath  string
}

func listFiles(response http.ResponseWriter, request *http.Request, username string) {
	files := make([]fileInfo, 0)

	//////////////////////////////////
	// BEGIN TASK 4: YOUR CODE HERE
	//////////////////////////////////

	// TODO: for each of the user's files, add a
	// corresponding fileInfo struct to the files slice.

	// replace this line
	fmt.Fprintf(response, "placeholder")

	//////////////////////////////////
	// END TASK 4: YOUR CODE HERE
	//////////////////////////////////

	data := map[string]interface{}{
		"Username": username,
		"Files":    files,
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/list.html")
	if err != nil {
		log.Error(err)
	}
	err = tmpl.Execute(response, data)
	if err != nil {
		log.Error(err)
	}
}

func getFile(response http.ResponseWriter, request *http.Request, username string) {
	fileString := strings.TrimPrefix(request.URL.Path, "/file/")

	_ = fileString

	//////////////////////////////////
	// BEGIN TASK 5: YOUR CODE HERE
	//////////////////////////////////

	// replace this line
	fmt.Fprintf(response, "placeholder")

	//////////////////////////////////
	// END TASK 5: YOUR CODE HERE
	//////////////////////////////////
}

func setNameOfServedFile(response http.ResponseWriter, fileName string) {
	response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
}

func processShare(response http.ResponseWriter, request *http.Request, sender string) {
	recipient := request.FormValue("username")
	filename := request.FormValue("filename")
	_ = filename

	if sender == recipient {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, "can't share with yourself")
		return
	}

	//////////////////////////////////
	// BEGIN TASK 6: YOUR CODE HERE
	//////////////////////////////////

	// replace this line
	fmt.Fprintf(response, "placeholder")

	//////////////////////////////////
	// END TASK 6: YOUR CODE HERE
	//////////////////////////////////

}

// Initiate a new session for the given username
func initSession(response http.ResponseWriter, username string) {
	// Generate session token
	sessionToken, err := randomByteString(16)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, err.Error())
		return
	}

	expires := time.Now().Add(sessionDuration)

	// Store session in database
	_, err = db.Exec("INSERT INTO sessions VALUES (NULL, ?, ?, ?)", username, sessionToken, expires.Unix())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, err.Error())
		return
	}

	// Set cookie with session data
	http.SetCookie(response, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expires,
		SameSite: http.SameSiteStrictMode,
	})
}
