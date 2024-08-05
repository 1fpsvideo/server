package main

import (
	"net/http"
	"os"
	"server/api/sessions"
	"server/api/system"
	"server/config"
	"server/debug"
	"server/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(config.UPLOAD_DIR, os.ModePerm); err != nil {
		panic(err)
	}

	// Route for the index page
	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	// Route for uploading the screenshot without authentication
	router.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	router.HandleFunc("/x/{sessionID}/ws", handlers.CursorWebSocketHandler)

	// Previously protected routes, now accessible without authentication
	router.HandleFunc("/x/{sessionID}/screenshot", handlers.ScreenshotHandler).Methods("GET")
	router.HandleFunc("/x/{sessionID}", handlers.SessionHandler).Methods("GET")

	// New API endpoints
	router.HandleFunc("/v1/api/system/version", system.VersionHandler).Methods("GET")
	router.HandleFunc("/v1/api/system/ping", system.PingHandler).Methods("GET")

	// Initialize ApiSessions
	apiSessions := sessions.NewApiSessions()
	router.HandleFunc("/v1/api/sessions", apiSessions.CreateSession).Methods("POST")

	debug.PrintDebug("Server is running on http://localhost:8899")
	http.ListenAndServe(":8899", router)
}
