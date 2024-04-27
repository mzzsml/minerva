package crtsh

import (
	"encoding/json"
	"io"
	"net/http"
)

func QueryCrtsh(w http.ResponseWriter, h *http.Request, d string) []string {
	var domainNames []string
	var m map[string]interface{}

	// Make the request to crt.sh, specifying the output type JSON
	resp, err := http.Get("https://crt.sh/?q=" + d + "&output=json")
	if err != nil {
		panic(err)
	}

	// resp.Body is byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Create an interface that will hold JSON structure
	// Unmarshal converts the JSON into the interface
	var f []interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		panic(err)
	}

	// Since the JSON returned is an array, we need to iterate through all JSONs
	domainNames = make([]string, len(f))
	for i, item := range f {
		m = item.(map[string]interface{})
		domainNames[i] = m["common_name"].(string)
	}

	var tempDomains = make([]string, 0, len(domainNames))
	tempDomains = append(tempDomains, domainNames[0])
	for _, d := range domainNames {
		for i := 0; i < len(tempDomains); i++ {
			if !isFound(d, tempDomains) {
				tempDomains = append(tempDomains, d)
			}
		}
	}

	return tempDomains
}

func isFound(target string, slice []string) bool {
	var found bool = false

	for _, item := range slice {
		if target == item {
			found = true
		}
	}

	return found
}
