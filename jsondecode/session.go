package jsondecode

import (
	"encoding/json"
	"strings"
)

type Sessions struct {
	Data struct {
		User []struct {
			Email      string `json:"email"`
			Username   string `json:"username"`
			ExternalId string `json:"external-id"`
		} `json:"user"`
		RememberToken string `json:"remember-token"`
		SessionToken  string `json:"session-token"`
	} `json:"data"`
	Context string `json:"context"`
}

func DecodeSessions(data string) (sessionToken string) {
	dec := json.NewDecoder(strings.NewReader(data))
	var sessions Sessions
	dec.Decode(&sessions)

	// fmt.Printf("Session-token: %s\n", sessions.Data.SessionToken)
	return sessions.Data.SessionToken

}
