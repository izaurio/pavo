package main

import (
	"encoding/json"
	"net/http"
)

// Get parameters for convert from Request query string
func GetConvertParams(req *http.Request) (map[string]string, error) {
	raw_converts := req.URL.Query().Get("converts")

	if raw_converts == "" {
		raw_converts = "{}"
	}

	convert := make(map[string]string)

	err := json.Unmarshal([]byte(raw_converts), &convert)
	if err != nil {
		return nil, err
	}

	return convert, nil
}
