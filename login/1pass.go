package login

import (
	"example/user/tasty/httpclient"
	"log"
	"os/exec"
	"strings"
)

// 1 Password secret references
const (
	op = "/usr/bin/op"

	sbxUserName   = "op://Private/Tasty_sbx/username"
	sbxPassword   = "op://Private/Tasty_sbx/credential"
	sbxVaultUser  = "op://Private/tastytrade-sbx-api/username"
	sbxVaultToken = "op://Private/tastytrade-sbx-api/credential"
	sbxApiItem    = "tastytrade-sbx-api"
)

func TrimNewLine(value string) (cleanString string) {

	cleanString = strings.TrimSuffix(value, "\n")

	return cleanString

}

func GetCreds(userRef string, passwordRef string) (string, string) {
	// get the credentials for the Tastytrade API stored in
	// 1Password Vault
	user_ref, err := exec.Command(op, "read", userRef).Output()

	if err != nil {
		log.Fatal("can't read secret reference for username ", err)
	}

	token_ref, err := exec.Command(op, "read", passwordRef).Output()

	if err != nil {
		log.Fatal("can't read secret reference for api token ", err)
	}

	username := string(user_ref[:])
	username = TrimNewLine(username)

	token := string(token_ref[:])
	token = TrimNewLine(token)

	return username, token
}

func WriteCreds(user string, sessionToken string) {
	// writes the session token to 1Password to be used for
	// API calls "op://Private/tastytrade-sbx-api/credential"

	// craft credential string
	credential := "credential=" + sessionToken
	// fmt.Printf("credential field: %s\n", credential)

	_, err := exec.Command(op, "item", "edit", sbxApiItem, credential).Output()

	if err != nil {
		log.Fatal("can't edit api session token credential ", err)
	}

}

func GetSessionToken(debug bool) {
	// login to get the session and remember tokens
	login, password := GetCreds(sbxUserName, sbxPassword)

	// trim the new line from the login value before returning
	login = strings.TrimSuffix(login, "\n")

	sessionToken := httpclient.GetSessionTokens(login, password, debug)

	// TODO debug
	// fmt.Printf("session token: %s", sessionToken)

	// write session token to 1Password
	WriteCreds(login, sessionToken)
}
