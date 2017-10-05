package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gh "github.com/rjz/githubhook"
)

var (
	secKey = os.Getenv("GH_SECRET")
	port   = os.Getenv("PORT")
)

func main() {
	if port == "" {
		port = "8080"
	}
	fmt.Printf("handler is listening on :8080/hook")

	http.HandleFunc("/hook", handleFunc)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Error(w http.ResponseWriter, f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+f+"\n", args)
	fmt.Fprintf(w, "Error: "+f, args)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleFunc(w http.ResponseWriter, req *http.Request) {
	hook, err := gh.Parse(nil, req)
	if err != nil {
		Error(w, "could not auth: %s", err)
		return
	}
	payload := json.RawMessage(hook.Payload)

	b, err := json.Marshal(struct {
		Payload *json.RawMessage
		Event   string
		Id      string
	}{
		Payload: &payload,
		Event:   hook.Event,
		Id:      hook.Id,
	})
	if err != nil {
		Error(w, "could not Marshal %s", err)
		return
	}

	_, err = os.Stdout.Write(b)
	if err != nil {
		Error(w, "while writing to stdout %s", err)
	}
}