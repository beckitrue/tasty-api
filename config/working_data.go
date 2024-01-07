package config

import (
	"os"
	"fmt"
	"bytes"
	"io"
	"encoding/json"
)

const (
	filename = "/home/becki/environment.json"
)

// Marshal object to io.Reader in JSON
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func GetWorkingData(v interface{}) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	
	return Unmarshal(jsonFile, v)
}

func WriteToJSON(v interface {}) error {
    jsonFile, err := os.Create(filename)
    if err != nil {
		return fmt.Errorf("error creating JSON file: %v", err)
	}

	r, err := Marshal(v)
	if err != nil {
		return err
	}

	_, err = io.Copy(jsonFile, r)

	// Close file and return any errors
	if err := jsonFile.Close(); err != nil {
		return fmt.Errorf("error closing JSON file: %v", err)
	}

	return err
}
