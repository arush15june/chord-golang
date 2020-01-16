package main

import (
	"fmt"
	"net/http"
)

// KeyLookupHandler in HTTP Handler for looking up keys in chord.
func KeyLookupHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		key := req.FormValue("key")
		host, err := LookupKey(key)
		if err != nil {
			fmt.Fprintf(w, "Lookup err: %v", err)
		}

		fmt.Fprintf(w, "%s\n", host)
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func InitHttpServer() {
	logger.Println("Initialized HTTP Server at port 8090")
	http.HandleFunc("/lookup", KeyLookupHandler)
	go http.ListenAndServe(":"+*ApiPort, nil)
}
