package crtsh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func QueryCrtsh(w http.ResponseWriter, h *http.Request, d string) {
	resp, err := http.Get("https://crt.sh/?q=" + d + "&output=json")
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var f []interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		panic(err)
	}

	for _, item := range f {
		m := item.(map[string]interface{})
		fmt.Fprintln(w, m["common_name"].(string))
	}
}
