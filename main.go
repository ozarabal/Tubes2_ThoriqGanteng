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

type response struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/submit", handleSubmit)
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
    response := response{Result: "OK", Message: "data : " + data.Start + ", " + data.Goal}
	log.Println("data : " + data.Start + ", " + data.Goal)
	fmt.Println(query.GetLinks("https://en.wikipedia.org/wiki/"+data.Start))
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}