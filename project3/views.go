// View-only functions for our web applications.
// Nothing here needs to change.
package main

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type PageData struct {
	Username, Error string
}

func NewPageData(username, error string) PageData {
	return PageData{Username: username, Error: error}
}

func showPage(response http.ResponseWriter, templateName string, data PageData) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+templateName+".html")
	if err != nil {
		log.Error(err)
	}
	err = tmpl.Execute(response, data)
	if err != nil {
		log.Error(err)
	}
}
