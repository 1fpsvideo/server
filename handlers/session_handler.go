package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"server/debug"

	"server/database"
	"server/handlers/templates"
	"server/repository"

	"github.com/gorilla/mux"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	debug.PrintDebug("Received index request")

	// Get the session ID from the URL parameters
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]

	debug.PrintDebug(fmt.Sprintf("Session ID: %s", sessionID))

	// Create a new RepoSessions instance
	repoSessions := repository.NewRepoSessions(database.GetRedisClient())

	// Check if the session exists in Redis
	sessionValue, err := repoSessions.Get(sessionID)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error getting session: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if sessionValue == nil {
		// Session not found, render the session_not_found template
		tmpl, err := template.ParseFS(templates.TemplateFS, "session_not_found.html")
		if err != nil {
			debug.PrintDebug(fmt.Sprintf("Error parsing template: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			debug.PrintDebug(fmt.Sprintf("Error executing template: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Session exists, render the share template
	tmpl, err := template.ParseFS(templates.TemplateFS, "x.html")
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error parsing template: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pass the session ID to the template
	data := struct {
		SessionID string
	}{
		SessionID: sessionID,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error executing template: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
