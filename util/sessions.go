package util

import "regexp"

// IsValidSessionID checks if the given sessionID is valid.
// It returns true if the sessionID consists of 1 to 15 alphanumeric characters.
func IsValidSessionID(sessionID string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]{1,15}$", sessionID)
	return match
}
