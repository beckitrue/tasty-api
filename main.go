package main

import (
	"example/user/tasty/httpclient"
	"example/user/tasty/jsondecode"
	"example/user/tasty/login"
	"fmt"
	"log"
)

const (
	sbxVaultUser  = "op://Private/tastytrade-sbx-api/username"
	sbxVaultToken = "op://Private/tastytrade-sbx-api/credential"
)

type ApiMsg struct {
	method string
	msg    string
	model  string
}

func main() {

	// // setting env to sbx for safety during testing
	env := "sbx"

	// login to TastyTrade to get the session token
	// and save it to 1Password for future API calls
	// https://developer.tastytrade.com/api-guides/sessions/
	// login.GetSessionToken()

	// get session token from 1Password
	username, token := login.GetCreds(sbxVaultUser, sbxVaultToken)
	fmt.Printf("Hello %s\ttoken:%s\n", username, token)

	// command we'll use for testing
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}

	cmd := accountList
	log.Printf("running command: %s", cmd)

	// configure the API URL
	url := httpclient.CreateURL(env, cmd.msg)

	// make the API call
	respString := httpclient.ApiCall(token, url, cmd.method)
	log.Println(respString)

	// perform the appropriate JSON decoding and
	// print output
	jsondecode.PrintDataAccounts(respString)

}
