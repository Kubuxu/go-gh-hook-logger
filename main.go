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
	file   = os.Getenv("STORAGE_FILE")
)

var f *os.File

func main() {
	if port == "" {
		port = "8080"
	}
	if secKey == "" {
		log.Fatal("secKey is empty")
	}
	fmt.Printf("handler is listening on :8080/hook")

	var err error
	f, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not open file '%s', err: %s\n", file, err)
		return
	}

	http.HandleFunc("/hook", hook)
	http.HandleFunc("/live", live)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Error(w http.ResponseWriter, f string, args ...interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Prefix()
	log.Printf("Error: "+f+"\n", args)
	fmt.Fprintf(w, "Error: "+f, args)
}

func live(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "It's ALIVE")
}

func hook(w http.ResponseWriter, req *http.Request) {
	hook, err := gh.Parse([]byte(secKey), req)
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

	_, err = f.Write(b)
	if err != nil {
		Error(w, "while writing to stdout %s", err)
	}
	_, err = f.Write([]byte{'\n'})
	if err != nil {
		Error(w, "while writing to stdout %s", err)
	}
}
