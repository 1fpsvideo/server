package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/database"
	"server/debug"
	"server/repository"
	"server/util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	debug.PrintDebug("Received upload request")

	sessionID := r.FormValue("session_id")
	if !util.IsValidSessionID(sessionID) {
		debug.PrintDebug("Invalid session ID format")
		http.Error(w, "Invalid session ID format", http.StatusBadRequest)
		return
	}

	repoSessions := repository.NewRepoSessions(database.GetRedisClient())
	session, err := repoSessions.Get(sessionID)
	if err != nil || session == nil {
		debug.PrintDebug("Session not found or expired")
		http.Error(w, "Session not found or expired. Please initiate a new session.", http.StatusNotFound)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error retrieving file from form: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create session-specific directory
	sessionDir := filepath.Join(config.UPLOAD_DIR, sessionID)
	err = os.MkdirAll(sessionDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		debug.PrintDebug(fmt.Sprintf("Error creating session directory %s: %v", sessionDir, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the uploaded file
	screenshotPath := filepath.Join(sessionDir, config.SCREENSHOT_FILE)
	out, err := os.Create(screenshotPath)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error creating file %s: %v", screenshotPath, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error writing file %s: %v", screenshotPath, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	debug.PrintDebug(fmt.Sprintf("Saved uploaded file to %s", screenshotPath))

	// Save the hash of the uploaded screenshot
	hash := calculateFileHash(screenshotPath)
	hashPath := filepath.Join(sessionDir, config.SCREENSHOT_HASH_FILE)
	err = os.WriteFile(hashPath, []byte(hash), 0644)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("Error writing hash file %s: %v", hashPath, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	debug.PrintDebug(fmt.Sprintf("Saved screenshot hash to %s", hashPath))

	w.WriteHeader(http.StatusOK)
}

func calculateFileHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
