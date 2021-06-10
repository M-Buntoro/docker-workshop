package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "RNG App")
	})

	http.HandleFunc("/rng", func(w http.ResponseWriter, r *http.Request) {

		maxStr := r.URL.Query().Get("max")
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			max = 25
		}

		rand.Seed(time.Now().UnixNano())
		randomized := rand.Intn(int(max))

		data := map[string]interface{}{
			"rand": randomized,
		}

		respBytes, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	})

	log.Println("Begin Serving RNG in port :11991")

	log.Println(http.ListenAndServe(":11991", nil))
}
