package query

import (
	"fmt"
	"io"
	// "io/ioutil"
	"net/http"
	// "regexp"
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

func GetPathIDS(start, goal string,allPaths *[][]string){
	startURL := "https://en.wikipedia.org/wiki/" + start
	goalURL := "https://en.wikipedia.org/wiki/" + goal
	fmt.Println("Start URL:", startURL)
	root := &Node{PageURL: startURL} // Membuat node root dengan startURL sebagai PageURL
	visited := make(map[string]bool) // Membuat map untuk menyimpan node yang sudah dieksplorasi

	// paths := [][]*Node{}
	paths, found := IDS(startURL, goalURL, visited, root)
	if !found{
		fmt.Println("shortest path not found")
		// return nil, fmt.Errorf("shortest path not found")
	}
	getAllPathIDS(paths,allPaths)
	// PrintAllPathIDS(*allPaths)
	
	// return paths, nil // Mengembalikan jalur terpendek dari startURL ke goalURL
}

func PrintAllPathIDS(paths [][]string){
	fmt.Println("Number of Paths:", len(paths))
	for i := range paths{
		fmt.Println("Path ke : ",i+1)
		for j := range paths[i] {
			fmt.Println(paths[i][j]) // Mencetak URL dari setiap node dalam jalur terpendek yang ditemukan
		}
	}
}

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

func DLS(limit int,goalURL string,mxLimit int, parent *Node) {
	defer wg.Done()
	cnt++
	fmt.Println("cnt : ",cnt)

	if time.Since(startT) > tOut {
		return
	}

	if parent.PageURL == goalURL { 
		pathsAns = append(pathsAns, getPath(parent))
		
	}

	if(limit <=0){
		return 
	}

	links, err := GetLinks(parent.PageURL)
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return 
	}
	var mxGo int
	if(mxLimit <=2){
		mxGo = 50
	}else if (mxLimit == 3){
		mxGo = 25
	}else if (mxLimit == 4){
		mxGo = 10
	}else{
		mxGo = 5
	}
	goCnt := 0
	for _, link := range links {
		child := &Node{PageURL: link, Parent: parent}
		
		parent.Children = append(parent.Children, child)
		currentLimit := limit-1
		wg.Add(1)

		var wg2 sync.WaitGroup
		wg2.Add(1)
		goCnt++
		
		go func(){
			defer wg2.Done()
			DLS(currentLimit, goalURL, mxLimit, child)
		}()
		if goCnt>=mxGo{
			wg2.Wait()
			goCnt = 0;
		}
			
	}
}

func IDS(startURL, goalURL string,visited map[string]bool, parent *Node) ([][]*Node, bool) {
	tOutSeconds = 290
	pathsAns = [][]*Node{}
	// fmt.Println(len(pathsAns))
	cnt = 0
	tOut = time.Duration(tOutSeconds) * time.Second
	startT = time.Now()
	for depth := 0; depth <= 6; depth++ {	
		if time.Since(startT) > tOut {
			if  len(pathsAns) != 0 {
				return pathsAns,true
			}else{
				return pathsAns,false
			}
		}
		fmt.Println("Depth:", depth)
		
		wg.Add(1)
		cnt ++
		DLS(depth,goalURL,depth,parent)
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
        cacheMux.RUnlock()
        // fmt.Println("Returning cached data")
        return doc, nil
    }
    cacheMux.RUnlock()

	client := &http.Client{}
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	links := []string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection)  {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, "Main_Page") && !strings.Contains(link, ":") {
			link = "https://en.wikipedia.org" + link
			links = append(links, link)
			// fmt.Println(link)
		}
	})	

	// Store in cache
    cacheMux.Lock()
    cache[pageURL] = links
    cacheMux.Unlock()
	
	return links, nil
}