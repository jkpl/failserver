package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func simulateLatency() {
	time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
}

func handler(w http.ResponseWriter, r *http.Request) {
	simulateLatency()

	switch n := rand.Intn(100); n {
	case 4:
		http.Error(w, "Could not find your lucky number!", http.StatusNotFound)
	case 5:
		http.Error(w, "Failed to compute your lucky number!", http.StatusInternalServerError)
	default:
		fmt.Fprintf(w, "Your lucky number is %d", n)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
