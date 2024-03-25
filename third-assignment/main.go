package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var (
	status     Status
	statusLock sync.Mutex
)

func updateStatus() {
	for {
		statusLock.Lock()
		status.Water = rand.Intn(100) + 1
		status.Wind = rand.Intn(100) + 1
		statusLock.Unlock()
		time.Sleep(15 * time.Second)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	statusLock.Lock()
	defer statusLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]Status{"status": status})
}

func getStatusHTML(w http.ResponseWriter, r *http.Request) {
	statusLock.Lock()
	defer statusLock.Unlock()

	waterStatus := ""
	if status.Water < 5 {
		waterStatus = "Aman"
	} else if status.Water >= 6 && status.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	windStatus := ""
	if status.Wind < 6 {
		windStatus = "Aman"
	} else if status.Wind >= 7 && status.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>Status</title>
        <meta http-equiv="refresh" content="2">
    </head>
    <body>
        <h1>Status</h1>
        <p>Water: %d meter - %s</p>
        <p>Wind: %d meter/second - %s</p>
    </body>
    </html>`, status.Water, waterStatus, status.Wind, windStatus)

	fmt.Fprint(w, html)
}

func main() {
	go updateStatus()

	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/", getStatusHTML)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started at http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
