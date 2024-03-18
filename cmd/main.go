package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/trace"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {

	fr := trace.NewFlightRecorder()
	if err := fr.Start(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler(fr))
	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(fr *trace.FlightRecorder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d := someSlowFunc()
		if d > 300 {
			var b bytes.Buffer
			_, err := fr.WriteTo(&b)
			if err != nil {
				log.Print(err)
				return
			}
			// Write it to a file.
			if err := os.WriteFile("trace.out", b.Bytes(), 0o755); err != nil {
				log.Print(err)
				return
			}
		}
		w.Write([]byte(fmt.Sprintf("slept for %d seconds", d)))
	}
}

func someSlowFunc() int {
	mi := 15
	ma := 3000
	delay := rand.Intn(ma-mi+1) + mi

	time.Sleep(time.Duration(delay) * time.Millisecond)
	return delay
}
