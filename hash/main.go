package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ComputeSha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hash App")
	})

	http.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {

		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		res := ComputeSha256(string(req))

		data := map[string]interface{}{
			"hash": res,
		}

		respBytes, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	})

	log.Println("Begin Serving RNG in port :11992")

	log.Println(http.ListenAndServe(":11992", nil))
}
