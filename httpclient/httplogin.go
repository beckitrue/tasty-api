package httpclient

import (
	"bytes"
	"encoding/json"
	"example/user/tasty/jsondecode"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetSessionTokens(login string, password string) (session string) {

	// HTTP endpoint
	sessionURL := "https://api.cert.tastyworks.com/sessions"

	// JSON body
	body := []byte(`{
		"login":"` + login + `",
		"password":"` + password + `"
	}`)

	// Create a HTTP post request
	req, err := http.NewRequest("POST", sessionURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// TODO: Debug
	// // debug http call
	// reqDump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("REQUEST:\n%s", string(reqDump))
	// // end debug

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
	}

	fmt.Print(string(resBody))

	sessionToken := jsondecode.DecodeSessions((string(resBody)))

	return sessionToken

}
