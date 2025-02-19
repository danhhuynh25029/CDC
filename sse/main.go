package main

import (
	"fmt"
	"net/http"
	"time"
)

var id = 0

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set http headers required for SSE
	fmt.Println("--------")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for client disconnection
	//clientGone := r.Context().Done()
	// set Timeout
	timeout := time.After(10 * time.Second)

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-timeout:
			fmt.Println("Client disconnected")
			return
		case <-t.C:
			// Send an event to the client
			// Here we send only the "data" field, but there are few others
			// data : contain data response client
			// retry : time to client reconnect to server when lost connection. You can reconnect with code in client
			// id : last-event-id
			_, err := fmt.Fprintf(w, "id: %d \ndata: The time is %s\nretry: 10000\n\n", id, time.Now().Format(time.UnixDate))
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
			id += 1
		}
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)
	fmt.Println("server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err.Error())
	}
}
