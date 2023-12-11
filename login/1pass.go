package login

import (
	"log"
	"os/exec"
	"strings"
)

// 1 Password secret references
const (
	op               = "/usr/bin/op"
	sbxVaultUser     = "op://Private/tastytrade-sbx-api/username"
	sbxVaultToken    = "op://Private/tastytrade-sbx-api/credential"
	sbxRememberToken = "op://Private/tastytrade-sbx-api/remember"
)

func GetCreds() (string, string) {
	// get the credentials for the Tastytrade API stored in
	// 1Password Vault
	user_ref, err := exec.Command(op, "read", sbxVaultUser).Output()

	if err != nil {
		log.Fatal("can't read secret reference for username ", err)
	}

	token_ref, err := exec.Command(op, "read", sbxVaultToken).Output()

	if err != nil {
		log.Fatal("can't read secret reference for api token ", err)
	}

	username := string(user_ref[:])
	token := string(token_ref[:])
	// trim the new line from the token value before returning
	token = strings.TrimSuffix(token, "\n")

	return username, token
}
