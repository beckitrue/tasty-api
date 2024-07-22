package httpclient

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"tastyapi/login"
)

func CreateURL(env string, endpoint string) (url string) {
	// creates the URL for the API endpoint that we want to call

	const (
		sbxURL  string = "https://api.cert.tastyworks.com/"
		prodURL string = "https://api.tastyworks.com/"
	)

	// default to sbx API endpoint for safety
	url = sbxURL

	if env == "prod" {
		url = prodURL
	}

	url += endpoint

	return url

}

func ApiCall(requestURL string, request string, debug bool) string {
    // set env for testing
    env := "sbx"

	// Get session or remember me token
	token := login.GetStoredToken(env)

	// TODO: logic for Method selection based on func input params
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
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

	// TODO: mask authorization token
	// debug http call
	if debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("REQUEST:\n%s", string(reqDump))
	}

	res, err := client.Do(req)

	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
	}

	if debug {
		log.Printf("client: status code: %d\n", res.StatusCode)
	}

	if res.StatusCode > 299 {
		log.Printf("Failed response client: status code: %d\n", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)

	// do some error checking on the response
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
	}

	if !json.Valid([]byte(resBody)) {
		log.Print("invalid JSON string returned: ", resBody)
	}

	return (string(resBody))

}
