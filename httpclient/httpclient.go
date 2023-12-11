package httpclient

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func CreateURL(env string, endpoint string) (url string) {
	// creates the URL for the API endpoint that we want to call

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

func ApiCall(token string, requestURL string, request string) string {

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
