package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/beckitrue/tasty-api/jsondecode"
)

const (
	sessionURL = "https://api.cert.tastyworks.com/sessions"
)

func DebugRequest(req *http.Request) {
	// print debug messages on the HTTP call

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("REQUEST:\n%s", string(reqDump))
}


func DestroySession(token string, debug bool) {

	req, err := http.NewRequest("DELETE", sessionURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authoriztion", token)
	req.Header.Add("User-Agent", "my-custom-client/2.0")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// debug http call
	if debug {
		DebugRequest(req)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed response client: error making http request: %s\n", err)
	}
	
	// log and quit if we get anything other than a 200 status code
	if res.StatusCode > 299 {
	    log.Fatalf("Failed response client status code: %d\n", res.StatusCode)
	}

}

func GetSessionTokens(login string, password string, debug bool) (session string) {

	// JSON body
	body := []byte(`{
		"login":"` + login + `",
		"password":"` + password + `"
	}`)

	// Create a HTTP post request
	req, err := http.NewRequest("POST", sessionURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "my-custom-client/2.0")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// debug http call
	if debug {
		DebugRequest(req)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed response client: error making http request: %s\n", err)
	}
    
	// log and quit if we get anything other than a 200 status code
	if res.StatusCode > 299 {
	    log.Fatalf("Failed response client: status code: %d\n", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)

	// do some error checking on the response
	if err != nil {
		log.Fatalf("client: could not read response body: %s\n", err)
	}

	if !json.Valid([]byte(resBody)) {
		log.Print("invalid JSON string returned: ", resBody)
	}

	fmt.Print(string(resBody))

	sessionToken := jsondecode.DecodeSessions((string(resBody)))

	return sessionToken

}
