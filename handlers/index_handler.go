package handlers

import (
	"html/template"
	"net/http"
	"server/debug"
	"server/handlers/templates"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	debug.PrintDebug("Received index request")

	// Parse the index template
	tmpl, err := template.ParseFS(templates.TemplateFS, "index.html")
	if err != nil {
		debug.PrintDebug("Error parsing index template: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		debug.PrintDebug("Error executing index template: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
