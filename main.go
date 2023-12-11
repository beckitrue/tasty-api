package main

import (
	"example/user/tasty/httpclient"
	"example/user/tasty/jsondecode"
	"example/user/tasty/login"
	"fmt"
	"log"
)

type ApiMsg struct {
	method string
	msg    string
	model  string
}

func main() {

	env := "sbx"

	// get session token from 1Password
	username, token := login.GetCreds()
	fmt.Printf("Hello %s here's your session token: %s\n", username, token)

	// command we'll use for testing
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts", model: "account"}

	cmd := accountList
	log.Printf("running command: %s", cmd)

	url := httpclient.CreateURL(env, cmd.msg)

	respString := httpclient.ApiCall(token, url, cmd.method)
	log.Println(respString)

	jsondecode.PrintDataAccounts(respString)

}
