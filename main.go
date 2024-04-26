package main 

import (
	"fmt"
	// "time"
	"Tubes2_ThoriqGanteng/query"
	"encoding/json"
	"net/http"
	"log"
)

type FormData struct {
    Start string `json:"start"`
    Goal  string `json:"goal"`
}

type Proses struct {
    Tipe string `json:"tipe"`
}

type Alignment struct {
    Pilihan string `json:"pilihan"`
}

type response struct {
    Result []string `json:"result"`
}


func main() {
    http.HandleFunc("/submit", handleSubmit)
    http.HandleFunc("/submitAlignment", handleSubmitAlignment)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	// Set header untuk CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    decoder := json.NewDecoder(r.Body)
    var data FormData
    err := decoder.Decode(&data)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    links, err := query.GetLinks("https://en.wikipedia.org/wiki/"+data.Start)
    if err != nil {
        log.Println("Error getting links:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    response := response{Result:links}
	log.Println("data : " + data.Start + ", " + data.Goal)
	// fmt.Println(response)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleSubmitAlignment(w http.ResponseWriter, r *http.Request) {
    // Set header untuk CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight OPTIONS request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Handle POST request
    decoder := json.NewDecoder(r.Body)
    var data Proses
    err := decoder.Decode(&data)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    pilihan := Alignment{Pilihan: data.Tipe}
    log.Println("data : " + data.Tipe)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(pilihan)
}