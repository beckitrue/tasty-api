package login

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/beckitrue/tasty-api/httpclient"
)

// 1 Password secret references
const (
	op = "/usr/bin/op" // path to op

	sbxUserName   = "op://SBX/Tasty_sbx/username"
	sbxPassword   = "op://SBX/Tasty_sbx/credential"
	sbxVaultUser  = "op://SBX/tastytrade-sbx-api/username"
	sbxVaultToken = "op://SBX/tastytrade-sbx-api/credential"
	// sbxRememberToken = "op://SBX/tastytrade-sbx-api/remember-token"
	sbxApiItem = "tastytrade-sbx-api"

	prodVaultUser  = "op://Private/Tasty-api/username"
	prodVaultToken = "op://Private/Tasty-api/credential"
)

func TrimNewLine(value string) (cleanString string) {

	cleanString = strings.TrimSuffix(value, "\n")

	return cleanString

}

func GetStoredToken(env string) string {
	// get the current session token stored in 1Password

	if env == "prod" {
		_, token := GetCreds(prodVaultUser, prodVaultToken)
		return token
	}

	// default to sbx for safety

	_, token := GetCreds(sbxVaultUser, sbxVaultToken)

	return token
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
	// API calls "op://SBX/tastytrade-sbx-api/credential"

	// craft credential string
	credential := "credential=" + sessionToken

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

	if debug {
		fmt.Printf("session token: %s", sessionToken)
	}

	// write session token to 1Password
	WriteCreds(login, sessionToken)
}

func DisableToken(debug bool) {
	// gets the current session token to pass to the API call to
	// delete the session token - takes no action on the 1Password item

	_, currentToken := GetCreds(sbxVaultUser, sbxVaultToken)

	httpclient.DestroySession(currentToken, debug)
}
