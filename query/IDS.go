package query

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Node adalah tipe data untuk node dalam struktur data graph
type Node struct {
	PageURL  string  // URL dari halaman Wikipedia
	Children []*Node // Anak-anak dari node
	Parent   *Node   // Parent dari node
}


// Fungsi untuk mendapatkan path menggunakan algoritma IDS
func GetPathIDS(start, goal string,allPaths *[][]string,method string){
	startURL := "https://en.wikipedia.org/wiki/" + start
	goalURL := "https://en.wikipedia.org/wiki/" + goal
	fmt.Println("Start URL:", startURL)
	root := &Node{PageURL: startURL} // Membuat node root dengan startURL sebagai PageURL

	paths, found := IDS(startURL, goalURL, root, method)

	if !found{
		fmt.Println("shortest path not found") // Kondisi jika path tidak ditemukan
	}
	if(method == "FIRST"){
		getFirstPathIDS(paths,allPaths)
	}else if(method == "ALL"){
		getAllPathIDS(paths,allPaths)
	}
}

func GetCnt () int {
	return cnt
}
// Fungsi untuk mecetak seluruh path
func PrintAllPathIDS(paths [][]string){
	fmt.Println("Number of Paths:", len(paths))
	for i := range paths{
		fmt.Println("Path ke : ",i+1)
		for j := range paths[i] {
			fmt.Println(paths[i][j]) // Mencetak URL dari setiap node dalam jalur terpendek yang ditemukan
		}
	}
}

// Fungsi untuk mendapatkan seluruh path
func getAllPathIDS(paths [][]*Node,allPaths *[][]string){
	for i := range paths{
		path := []string{}
		for _, node := range paths[i] {
			prefix := "https://en.wikipedia.org/wiki/"
			cleanURL := strings.TrimPrefix(node.PageURL, prefix)
			path = append(path, cleanURL) 
		}
		*allPaths = append(*allPaths,path)
	}
}

// Fungsi untuk mendapatkan path pertama
func getFirstPathIDS(paths [][]*Node,allPaths *[][]string){
	path := []string{}
	for _, node := range paths[0] {
		prefix := "https://en.wikipedia.org/wiki/"
		cleanURL := strings.TrimPrefix(node.PageURL, prefix)
		path = append(path, cleanURL) 
	}
	*allPaths = append(*allPaths,path)
}

// Fungsi untuk mendapatkan path dari node leaf
func getPath(leaf *Node) []*Node {
	path := []*Node{leaf} // Mulai dengan node leaf sebagai bagian dari jalur
	current := leaf

	for current.Parent != nil { // Selama node saat ini memiliki parent
		parent := current.Parent // Dapatkan parent dari node saat ini
		path = append([]*Node{parent}, path...) // Masukkan parent ke dalam jalur di depan
		current = parent // Pindah ke parent node
	}

	return path
}



var wg sync.WaitGroup
var pathsAns = [][]*Node{}
var cnt int

var tOutSeconds int
var tOut time.Duration
var startT time.Time

// Fungsi DLS
func DLS(limit int,goalURL string,mxLimit int, parent *Node,method string) {
	defer wg.Done() 
	cnt++ // Mengitung jumlah link yang di cek
	// fmt.Println("cnt : ",cnt)

	if (method == "FIRST"){
		if  len(pathsAns) != 0 { // Jika method adalah FIRST dan path ditemukan, keluar dari fungsi
			return 
		}
	}

	if time.Since(startT) > tOut {
		return // Jika waktu habis keluar dari fungsi
	}

	if parent.PageURL == goalURL { 
		pathsAns = append(pathsAns, getPath(parent)) // Jika goalURL ditemukan tambahkan ke pathsAns
		if (method == "FIRST"){
			return
		}
	}

	if(limit <=0){
		return // Jika sudah melebihi depth keluar dari fungsi
	}

	links, err := GetLinks(parent.PageURL) // Ambil seluruh link dari page url
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return 
	}
	var mxGo int
	if(mxLimit <=2){
		mxGo = 50  // Depth <= 2 maksimal goroutine adalah 50 
	}else if (mxLimit == 3){
		mxGo = 25	// Depth = 3 maksimal goroutine adalah 25
	}else if (mxLimit == 4){
		mxGo = 10 // Depth = 4 maksimal goroutine adalah 10
	}else{
		mxGo = 5 // Depth > 4 maksimal goroutine adalah 50
	}
	goCnt := 0 // counter untuk menghitung banyak goroutine yang sedang berjalan

	for _, link := range links { // Iterasi seluruh link yang didapatkan
		child := &Node{PageURL: link, Parent: parent}
		parent.Children = append(parent.Children, child) // Tambahkan Node link ke parent 
		currentLimit := limit-1
		wg.Add(1)

		var wg2 sync.WaitGroup
		wg2.Add(1)
		goCnt++
		
		go func(){
			defer wg2.Done()
			DLS(currentLimit, goalURL, mxLimit, child,method) // Panggil DLS untuk kedalaman selanjutnya
		}()
		if goCnt>=mxGo{
			wg2.Wait()
			goCnt = 0;
		}
		if (method == "FIRST"){
			if  len(pathsAns) != 0 {
				return 
			}
		}
		if time.Since(startT) > tOut {
			return // Jika waktu habis keluar dari fungsi
		}
	}
}

// Fungsi IDS
func IDS(startURL, goalURL string, parent *Node, method string) ([][]*Node, bool) {
	tOutSeconds = 290 // maksimum time out 290 detik
	pathsAns = [][]*Node{}
	cnt = 0
	tOut = time.Duration(tOutSeconds) * time.Second
	startT = time.Now()
	for depth := 0; depth <= 6; depth++ {	// Iterasi depth mulai dari 0 sampai 6
		if time.Since(startT) > tOut {
			if  len(pathsAns) != 0 {
				return pathsAns,true // Jika waktu sudah habis kembalikan pathsAns
			}else{
				return pathsAns,false
			}
		}
		fmt.Println("Depth:", depth)
		
		wg.Add(1)
		cnt ++
		DLS(depth,goalURL,depth,parent,method) // Panggil DLS dengan parent adalah root
		wg.Wait()
		if  len(pathsAns) != 0 {
			
			return pathsAns,true
		}
	}
	return pathsAns,false
}


var (
    cache    = make(map[string][]string)
    cacheMux = sync.RWMutex{}
)

// getLinks mengambil tautan-tautan dari halaman Wikipedia
func GetLinks(pageURL string) ([]string, error) {
	cacheMux.RLock()
    if doc, found := cache[pageURL]; found {
        cacheMux.RUnlock() // Cek apakah pageURL telah tersimpan di cache
        return doc, nil
    }
    cacheMux.RUnlock()

	// Membuat HTTP request untuk halaman URL yang diberikan
	client := &http.Client{}
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// Mengirimkan request dan menerima response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Memeriksa status response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	// Membaca isi dari response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Membuat dokumen HTML dari isi body
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	cekLink := make(map[string]bool)
	cekLink [pageURL] = true

	// Menemukan dan mengekstrak tautan-tautan yang valid dalam dokumen
	links := []string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection)  {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, "Main_Page") && !strings.Contains(link, ":") {
			link = "https://en.wikipedia.org" + link
			if !cekLink[link]{ // Cek link unik atau tidak
				links = append(links, link) 
				cekLink[link] = true
			}
			
		}
	})	

	// Menyimpan hasil tautan di cache
    cacheMux.Lock()
    cache[pageURL] = links
    cacheMux.Unlock()
	
	return links, nil
}