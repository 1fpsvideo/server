package sessions

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"server/database"
	"server/repository"
)

const (
	SESSION_ID_LENGTH = 10
	SESSION_EXPIRY    = 48 * time.Hour
)

type ApiSessions struct {
	repoSessions *repository.RepoSessions
}

func NewApiSessions() *ApiSessions {
	return &ApiSessions{
		repoSessions: repository.NewRepoSessions(database.GetRedisClient()),
	}
}

func (api *ApiSessions) CreateSession(w http.ResponseWriter, r *http.Request) {
	sessionID := generateRandomID(SESSION_ID_LENGTH)

	err := api.repoSessions.Create(sessionID, SESSION_EXPIRY)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "ok",
		"session_id": sessionID,
	})
}

func generateRandomID(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
