package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		info := Info{
			Name:    "MasjidPi",
			Version: "0.1.0",
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(info)
	})

	log.Println("Starting MasjidPi on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
