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

	// Protected routes
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(authMiddleware)
	protectedRouter.HandleFunc("/x/{sessionID}/screenshot", handlers.ScreenshotHandler).Methods("GET")
	protectedRouter.HandleFunc("/x/{sessionID}", handlers.SessionHandler).Methods("GET")

	// New API endpoints
	router.HandleFunc("/v1/api/system/version", system.VersionHandler).Methods("GET")
	router.HandleFunc("/v1/api/system/ping", system.PingHandler).Methods("GET")

	// Initialize ApiSessions
	apiSessions := sessions.NewApiSessions()
	router.HandleFunc("/v1/api/sessions", apiSessions.CreateSession).Methods("POST")

	debug.PrintDebug("Server is running on http://localhost:4567")
	http.ListenAndServe(":4567", router)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "sky" || password != "31337" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted Area"`)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			debug.PrintDebug("Authentication failed")
			return
		}
		debug.PrintDebug("Authentication successful")
		next.ServeHTTP(w, r)
	})
}
