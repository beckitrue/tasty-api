package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	username := os.Getenv("TT_USER")
	session_token := os.Getenv("TT_SESSION")

	requestURL := "https://api.cert.tastyworks.com/customers/me"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "my-custom-client/2.0")
	req.Header.Add("Authorization", session_token)

	// need to ignore cert in the Tastytrade sandbox
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: &config}
	client := &http.Client{Transport: tr}

	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

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
