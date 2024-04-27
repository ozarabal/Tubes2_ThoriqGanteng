package main

import (
	"fmt"

	"time"
	"Tubes2_ThoriqGanteng/query"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
    "github.com/rs/cors"
	"net/url"
    "strconv"
)

var mode string = ""
var pilihan_tipe string = ""

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
	Result [][]string `json:"result"`
    Waktu string `json:"waktu"`
    Banyak_path string `json:"banyak_path"`
    Banyak_jelajah string `json:"banyak_jelajah"`
    Kedalaman string `json:"kedalaman"`

}

func main() {
	router := mux.NewRouter()

    // Mengaktifkan CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type"},
        Debug:          true, 
    })
    handler := c.Handler(router)

	router.HandleFunc("/submit", handleSubmit).Methods("POST")
	router.HandleFunc("/submitAlignment", handleSubmitAlignment).Methods("POST")
	router.HandleFunc("/fetch-wikipedia", fetchWikipediaHandler).Methods("GET")
	router.HandleFunc("/submitmethod", handleSubmitmethod).Methods("POST")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var data FormData
	err := decoder.Decode(&data)
    var boole bool
    if (pilihan_tipe == "") {
        log.Println("Error: method not set")
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }else if (pilihan_tipe == "First") {
        boole = true
    }else if (pilihan_tipe == "All") {
        boole = false
    }
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
    if mode == "" {
        log.Println("Error: mode not set")
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    } else if mode == "BFS" {
        var urls []string
        graph := query.NewGraph()
        path := make(map[string]bool)
        path["https://en.wikipedia.org/wiki/"+data.Start] = true
        urls = append(urls, "https://en.wikipedia.org/wiki/"+data.Start) 
		time_start := time.Now()
        graph, path = query.Bfs2(urls, path, graph,"https://en.wikipedia.org/wiki/"+data.Start, "https://en.wikipedia.org/wiki/"+data.Goal,boole)
        
        visitied := make(map[string]bool)
        temppath := []string{}
        allpath := [][]string{}
        

        query.GetAllPaths(graph, "https://en.wikipedia.org/wiki/"+data.Start, "https://en.wikipedia.org/wiki/"+data.Goal, visitied, temppath, &allpath)
		time_end := time.Now()
		elapsed := time_end.Sub(time_start)
        
        elapsedSeconds := int(elapsed.Seconds())
        banyakpath := len(allpath)
        banyaklink := len(path)
        kedalamanpath := len(allpath[0])
        strBanyakPath := strconv.Itoa(banyakpath)
        strBanyakLink := strconv.Itoa(banyaklink)
        strKedalamanPath := strconv.Itoa(kedalamanpath)
        strwaktu := strconv.Itoa(elapsedSeconds)

        response  := response{Result:allpath, Waktu:strwaktu, Banyak_path:strBanyakPath, Banyak_jelajah:strBanyakLink, Kedalaman: strKedalamanPath}
        log.Println("mode : " + mode)
        log.Println("data : " + data.Start + ", " + data.Goal)
        // fmt.Println(response)
    
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }else if mode ==  "IDS" {
		fmt.Println(data.Start,data.Goal)
		allpath := [][]string{}
        time_start := time.Now()
        if pilihan_tipe == "First" {   
		    query.GetPathIDS(data.Start,data.Goal,&allpath,"FIRST")
        } else if pilihan_tipe == "All" {
            query.GetPathIDS(data.Start,data.Goal,&allpath,"ALL")
        }
		query.PrintAllPathIDS(allpath)
        time_end := time.Now()
        elapsed := time_end.Sub(time_start)
        banyakpath := len(allpath)
        banyaklink := query.GetCnt()
        kedalamanpath := len(allpath[0])
        strBanyakPath := strconv.Itoa(banyakpath)
        strBanyakLink := strconv.Itoa(banyaklink)
        strKedalamanPath := strconv.Itoa(kedalamanpath)
        strwaktu := strconv.Itoa(int(elapsed.Seconds()))


        response  := response{Result:allpath, Waktu:strwaktu, Banyak_path:strBanyakPath, Banyak_jelajah:strBanyakLink, Kedalaman: strKedalamanPath}
        log.Println("mode : " + mode)
        log.Println("data : " + data.Start + ", " + data.Goal)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
	}//else if mode == "IDS" {
    //     links, err := query.GetLinks("https://en.wikipedia.org/wiki/"+data.Start)
    //     if err != nil {
    //         log.Println("Error getting links:", err)
    //         http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    //         return
    //     }
    //     log.Println("IDS")
    //     response := response{Result:links}
    //     log.Println("mode : " + mode)
    //     log.Println("data : " + data.Start + ", " + data.Goal)
    //     // fmt.Println(response)
        
    //     w.Header().Set("Content-Type", "application/json")
    //     json.NewEncoder(w).Encode(response)
    // }
    
}

func handleSubmitAlignment(w http.ResponseWriter, r *http.Request) {
	
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
	if data.Tipe == "BFS" {
		mode = "BFS"
	} else if data.Tipe == "IDS" {
		mode = "IDS"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilihan)
}

func handleSubmitmethod(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body)
	var data Proses
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	pilihan := Alignment{Pilihan: data.Tipe}
	log.Println("pilihan : " + data.Tipe)
	if data.Tipe == "First" {
		pilihan_tipe = "First"
	} else if data.Tipe == "All" {
		pilihan_tipe = "All"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilihan)
}

const wikipediaAPIURL = "https://en.wikipedia.org/w/api.php"

func fetchWikipediaHandler(w http.ResponseWriter, r *http.Request) {
    searchQuery := r.URL.Query().Get("search")
    limit := "5"

    // Membuat permintaan ke API Wikipedia
    reqURL := wikipediaAPIURL + "?action=opensearch&search=" + url.QueryEscape(searchQuery) + "&limit=" + limit + "&namespace=0&format=json"
    resp, err := http.Get(reqURL)
    if err != nil {
        log.Println("Error fetching data from Wikipedia API:", err)
        http.Error(w, "Error fetching data from Wikipedia API", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Membaca dan mengirimkan kembali respons dari API Wikipedia
    var responseData interface{}
    err = json.NewDecoder(resp.Body).Decode(&responseData)
    if err != nil {
        log.Println("Error decoding JSON response from Wikipedia API:", err)
        http.Error(w, "Error decoding JSON response from Wikipedia API", http.StatusInternalServerError)
        return
    }

    // Mengatur header CORS untuk mengizinkan akses dari domain klien
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Mengirimkan kembali respons JSON ke klien
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}