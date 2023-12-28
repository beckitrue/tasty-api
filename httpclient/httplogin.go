package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/beckitrue/tasty-api/jsondecode"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	sessionURL = "https://api.cert.tastyworks.com/sessions"
)

func DestroySession (token string, debug bool) {

	req, err := http.NewRequest("DELETE", sessionURL, nil)
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authoriztion", token)

	if err != nil {
        panic(err)
    }

	client := &http.Client{}
 
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
	    log.Fatalf("client: error making http request: %s\n", err)     
	}                                                                  
	                                                                   
	log.Printf("client: status code: %d\n", res.StatusCode)

}

func GetSessionTokens(login string, password string, debug bool) (session string) {

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
