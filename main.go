package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const port = ":8080"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("random num: " + strconv.Itoa(rand.Int())))
	})
	log.Printf("server is starting on port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("got err: %v", err)
	}
}
