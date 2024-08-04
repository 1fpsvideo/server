package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/debug"
	"server/util"

	"github.com/gorilla/mux"
)

func ScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	debug.PrintDebug("Received screenshot request")

	// Extract sessionID from the URL parameters
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]

	if !util.IsValidSessionID(sessionID) {
		debug.PrintDebug("Invalid session ID format")
		http.Error(w, "Invalid session ID format", http.StatusBadRequest)
		return
	}

	screenshotPath := filepath.Join(config.UPLOAD_DIR, sessionID, config.SCREENSHOT_FILE)
	if _, err := os.Stat(screenshotPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	hashPath := filepath.Join(config.UPLOAD_DIR, sessionID, config.SCREENSHOT_HASH_FILE)
	hash, err := os.ReadFile(hashPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("ETag", string(hash))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, screenshotPath)

	debug.PrintDebug(fmt.Sprintf("Served screenshot file %s", screenshotPath))
}
