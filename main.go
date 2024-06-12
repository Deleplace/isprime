package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		if user == "" {
			http.Error(w, "mandatory parameter: user", http.StatusBadRequest)
			return
		}
		number := r.FormValue("number")
		if number == "" {
			http.Error(w, "mandatory parameter: number", http.StatusBadRequest)
			return
		}
		n := new(big.Int)
		_, ok := n.SetString(number, 10)
		if !ok {
			http.Error(w, "could not parse number \""+number+"\"", http.StatusBadRequest)
			return
		}
		slog.Info("will test primality",
			"user", user,
			"number", number)

		start := time.Now()
		isprime := n.ProbablyPrime(10_000_000)
		duration := time.Since(start)

		slog.Info("tested primality",
			"user", user,
			"number", number,
			"isPrime", isprime,
			"durationNs", duration.Nanoseconds(),
		)
		fmt.Fprintln(w, "is", n, "probably prime:", isprime)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	addr := os.Getenv("ADDR") + ":" + port
	log.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}
