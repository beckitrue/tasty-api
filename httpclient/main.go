package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

// API structs
type apiMsg struct {
	method string
	msg    string
}

func getCreds() (string, string) {
	// get the credentials for the Tastytrade API stored in
	// 1Password Vault

	// set the values for 1Password CLI secret references
	const (
		op            = "/usr/bin/op"
		sbxVaultUser  = "op://Private/tastytrade-sbx-api/username"
		sbxVaultToken = "op://Private/tastytrade-sbx-api/credential"
	)
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

	return username, token
}

func createURL(token string, env string, message apiMsg) {

	const (
		sbxURL  string = "https://api.cert.tastyworks.com/"
		prodURL string = "https://api.tastyworks.com/"
	)

	if env == "sbx" {
		requestURL := sbxURL + message.msg
		fmt.Println(requestURL)
		apiCall(token, requestURL, "GET")
	}

}

func apiCall(token string, requestURL string, request string) {

	// // requestURL := "https://api.cert.tastyworks.com/customers/me"
	// req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	fmt.Printf("token at client: %s\n", token)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "my-custom-client/2.0")
	req.Header.Add("Authorization", token)

	// need to ignore cert in the Tastytrade sandbox
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: &config}
	client := &http.Client{Transport: tr}

	// // client := &http.Client{
	// // 	Timeout: 10 * time.Second,
	// // }

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("client: error making http request: %s\n", err)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("client: could not read response body: %s\n", err)
	}
	fmt.Printf("client: response body: %s\n", resBody)

}

func main() {

	const (
		env = "sbx"
	)

	// get session token
	username, token := getCreds()

	fmt.Printf("Hello %s here's your session token: %s\n", username, token)

	// get my account info
	accountInfo := apiMsg{method: "GET", msg: "customer/me"}

	createURL(token, env, accountInfo)

}
