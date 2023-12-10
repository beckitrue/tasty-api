package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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

type ApiMsg struct {
	method string
	msg    string
}

type Context struct {
	context string
}

type Accounts struct {
	Data struct {
		Items []struct {
			Account struct {
				AccountNumber        string `json:"account-number"`
				OpenedAt             string `json:"opened-at"`
				Nickname             string `json:"nickname"`
				AccountTypeName      string `json:"account-type-name"`
				DayTraderStatus      bool   `json:"day-trader-status"`
				IsClosed             bool   `json:"is-closed"`
				IsFirmError          bool   `json:"is-firm-error"`
				IsFirmProprietary    bool   `json:"is-firm-proprietary"`
				IsFuturesApproved    bool   `json:"is-futures-approved"`
				IsTestDrive          bool   `json:"is-test-drive"`
				MarginOrCash         string `json:"margin-or-cash"`
				IsForeign            bool   `json:"is-foreign"`
				InvestmentObjective  string `json:"investment-objective"`
				SuitableOptionsLevel string `json:"suitable-options-level"`
				CreatedAt            string `json:"created-at"`
			} `json:"account"`
			AuthorityLevel string `json:"authority-level"`
		} `json:"items"`
	} `json:"data"`
	Context string `json:"context"`
}

func getCreds() (string, string) {
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

func createURL(env string, endpoint string) (url string) {

	const (
		sbxURL  string = "https://api.cert.tastyworks.com/"
		prodURL string = "https://api.tastyworks.com/"
	)

	// default to sbx API endpoint for safety
	baseURL := sbxURL

	if env == "prod" {
		baseURL = prodURL
	}

	baseURL += endpoint

	return baseURL

}

func apiCall(token string, requestURL string, request string) string {

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}

	// set header values
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "my-custom-client/2.0")
	req.Header.Add("Authorization", token)

	// need to ignore cert in the Tastytrade sandbox
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: &config}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * tr.IdleConnTimeout,
	}

	// debug http call
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("REQUEST:\n%s", string(reqDump))
	// end debug

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("client: error making http request: %s\n", err)
	}

	log.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)

	// do some error checking on the response
	if err != nil {
		log.Fatalf("client: could not read response body: %s\n", err)
	}

	if !json.Valid([]byte(resBody)) {
		log.Print("invalid JSON string returned: ", resBody)
		return (string(resBody))
	}

	return (string(resBody))

}

func getJson(data string) {

	dec := json.NewDecoder(strings.NewReader(data))
	for {
		var accounts Accounts
		if err := dec.Decode(&accounts); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// loop through the number items
		for i := range accounts.Data.Items {
			fmt.Printf("Account %d: %s\n", i, accounts.Data.Items[i].Account.AccountNumber)
		}
	}
}

func main() {

	// setting env to sbx for safety while developing and testing
	const (
		env = "sbx"
	)

	// get session token from 1Password
	username, token := getCreds()
	fmt.Printf("Hello %s here's your session token: %s\n", username, token)

	// command list - we'll use cli args eventually
	// accountInfo := apiMsg{method: "GET", msg: "customers/me"}
	accountList := ApiMsg{method: "GET", msg: "customers/me/accounts"}

	cmd := accountList
	log.Printf("running command: %s", cmd)

	url := createURL(env, cmd.msg)

	resString := apiCall(token, url, cmd.method)
	fmt.Println(resString)

	getJson(resString)

}
