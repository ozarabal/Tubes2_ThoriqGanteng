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

func FindTree(start, goal string, maxdepth int) ([]*Node, error) {
	startURL := "https://en.wikipedia.org/wiki/" + start
	goalURL := "https://en.wikipedia.org/wiki/" + goal
	fmt.Println("Start URL:", startURL)
	paths, err := ShortestPath(startURL, goalURL, maxdepth)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	fmt.Println("Number of Paths:", len(paths))
	fmt.Println("Shortest Path:")
	for _, node := range paths[0] {
		fmt.Println(node.PageURL) // Mencetak URL dari setiap node dalam jalur terpendek yang ditemukan
	}

	return paths[0], nil // Mengembalikan jalur terpendek dari startURL ke goalURL
}

func ShortestPath(startURL, goalURL string, maxDepth int) ([][]*Node, error) {
	root := &Node{PageURL: startURL} // Membuat node root dengan startURL sebagai PageURL
	visited := make(map[string]bool) // Membuat map untuk menyimpan node yang sudah dieksplorasi

	paths := [][]*Node{} // Menyimpan semua jalur yang terbentuk
	paths, found := IDS(startURL, goalURL, visited, root, paths)
	// IDS(startURL, goalURL string,visited map[string]bool, parent *Node, paths [][]*Node)

	if !found {
		return nil, fmt.Errorf("shortest path not found")
	}
	// if len(paths) == 0 {
	// 	return nil, fmt.Errorf("shortest path not found")
	// }

	return paths, nil
}

func getPath(leaf *Node) []*Node {
	path := []*Node{leaf} // Mulai dengan node leaf sebagai bagian dari jalur
	current := leaf

	for current.Parent != nil { // Selama node saat ini memiliki parent
		parent := current.Parent                // Dapatkan parent dari node saat ini
		path = append([]*Node{parent}, path...) // Masukkan parent ke dalam jalur di depan
		current = parent                        // Pindah ke parent node
	}

	return path
}

// getPaths mengambil jalur-jalur yang dapat dicapai dari startURL ke goalURL
func getPaths(startURL, goalURL string, visited map[string]bool, parent *Node, paths [][]*Node, depth, maxDepth int, found *bool) [][]*Node {
	if depth > maxDepth || *found {
		return paths // Melewati batas kedalaman maksimum atau jalur sudah ditemukan, berhenti pencarian untuk jalur ini
	}

	if parent.PageURL == goalURL { // Jika node saat ini adalah goal node
		// Path ditemukan, tambahkan jalur ke paths
		*found = true // Set penanda bahwa jalur telah ditemukan
		paths = append(paths, getPath(parent))
		return paths
	}

	if visited[parent.PageURL] { // Jika node sudah dieksplorasi sebelumnya
		return paths
	}
	visited[parent.PageURL] = true // Tandai node sebagai sudah dieksplorasi

	// Ambil tautan-tautan dari halaman node menggunakan GetLinks
	links, err := GetLinks(parent.PageURL)
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return paths
	}

	// Untuk setiap tautan, buat node anak dan rekursif cari jalur dengan penambahan kedalaman
	for _, link := range links {
		fmt.Println(link)
		child := &Node{PageURL: link, Parent: parent}
		parent.Children = append(parent.Children, child)
		paths = getPaths(startURL, goalURL, visited, child, paths, depth+1, maxDepth, found)
		if *found {
			break // Jika jalur sudah ditemukan, hentikan pencarian lebih lanjut
		}
	}

	return paths
}

var wg sync.WaitGroup
var pathsAns = [][]*Node{}
var cnt = 0

func DLS(limit int, goalURL string, visited map[string]bool, parent *Node, paths [][]*Node) {
	defer wg.Done()
	if parent.PageURL == goalURL {
		paths = append(paths, getPath(parent))
		pathsAns = make([][]*Node, len(paths))
		copy(pathsAns, paths)
		return
	}

	if limit <= 0 {
		return
	}

	links, err := GetLinks(parent.PageURL)
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return
	}
	for _, link := range links {
		child := &Node{PageURL: link, Parent: parent}

		if !visited[child.PageURL] {
			visited[child.PageURL] = true
			parent.Children = append(parent.Children, child)
			currentLimit := limit - 1
			wg.Add(1)

			var wg2 sync.WaitGroup
			wg2.Add(1)
			cnt++
			fmt.Println("cnt : ", cnt)
			go func() {
				defer wg2.Done()
				DLS(currentLimit, goalURL, visited, child, paths)
			}()
			wg2.Wait()
			if len(pathsAns) != 0 {

				return
			}
		}
	}
}

func IDS(startURL, goalURL string, visited map[string]bool, parent *Node, paths [][]*Node) ([][]*Node, bool) {

	for depth := 0; depth <= 6; depth++ {
		fmt.Println("Depth:", depth)
		visited := make(map[string]bool)
		wg.Add(1)
		cnt++
		DLS(depth, goalURL, visited, parent, paths)
		wg.Wait()
		if len(pathsAns) != 0 {

			return pathsAns, true
		}
	}
	return pathsAns, false
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
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") {
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

type Graph struct {
	adjacencyList map[string][]string
	visited       map[string]bool
}

func NewGraph() *Graph {
	return &Graph{
		adjacencyList: make(map[string][]string),
		visited:       make(map[string]bool),
	}
}

func (g *Graph) AddEdge(src, dest string) {
	g.adjacencyList[src] = append(g.adjacencyList[src], dest)
}

func PrintGraph(graph *Graph, start, end string) {
	visited := make(map[string]bool)
	path := []string{start}
	var dfs func(node string) bool
	dfs = func(node string) bool {
		if node == end {
			fmt.Println("Path:", strings.Join(path, " -> "))
			return true
		}
		visited[node] = true
		for _, neighbor := range graph.adjacencyList[node] {
			if !visited[neighbor] {
				path = append(path, neighbor)
				if dfs(neighbor) {
					return true
				}
				path = path[:len(path)-1]
			}
		}
		return false
	}
	dfs(start)
}

func (g *Graph) maxDepth(node string) int {
	g.visited = make(map[string]bool) // Reset visited map for each call
	return g.searchMax(node)
}

func (g *Graph) searchMax(node string) int {
	g.visited[node] = true
	maxDepth := 0
	for _, neighbor := range g.adjacencyList[node] {
		if !g.visited[neighbor] {
			depth := g.searchMax(neighbor)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}
	return 1 + maxDepth
}

func PrintAllPaths(graph *Graph, start, end string, visited map[string]bool, path []string) {
	visited[start] = true
	path = append(path, start)
	if start == end {
		fmt.Println("Path:", strings.Join(path, " -> "))
	} else {
		for _, neighbor := range graph.adjacencyList[start] {
			if !visited[neighbor] {
				PrintAllPaths(graph, neighbor, end, visited, path)
			}
		}
	}
	path = path[:len(path)-1]
	visited[start] = false
}

func (g *Graph) PrintAllPaths2() {
	for src := range g.adjacencyList {
		g.visited = make(map[string]bool)
		g.printAllPathsUtil(src, src, []string{})
	}
}

func (g *Graph) printAllPathsUtil(start, current string, path []string) {
	g.visited[current] = true
	path = append(path, current)

	if start != current {
		fmt.Println(path)
	}

	for _, v := range g.adjacencyList[current] {
		if !g.visited[v] {
			g.printAllPathsUtil(start, v, path)
		}
	}

	// Mark the current node as unvisited to explore other paths
	g.visited[current] = false
}

func GetAllPaths(graph *Graph, start, end string, visited map[string]bool, path []string, allPaths *[][]string) {
	visited[start] = true
	path = append(path, start)
	if start == end {
		*allPaths = append(*allPaths, append([]string{}, path...))
	} else {
		for _, neighbor := range graph.adjacencyList[start] {
			if !visited[neighbor] {
				GetAllPaths(graph, neighbor, end, visited, path, allPaths)
			}
		}
	}
	path = path[:len(path)-1]
	visited[start] = false
}

func PrintAllPaths2(allPaths [][]string) {
	for _, path := range allPaths {
		fmt.Println("Path:", strings.Join(path, " -> "))
	}
}

func Bfs(links []string, query map[string]bool, graph *Graph, final string, depth int) *Graph {
	found := false
	fmt.Println(depth)
	fmt.Println(len(links))
	fmt.Println(len(query))
	i := 0
	var newLinks []string
	for _, link := range links {
		i++
		//fmt.Println(i)
		links2, query2, graph2, found2 := getLinks(link, query, graph, final)
		query = query2
		newLinks = append(newLinks, links2...)
		graph = graph2
		if found2 == true {
			found = true
		}
		if found2 == true {
			fmt.Println(len(query2))
			break
		}
	}
	//graph.PrintAllPaths2()
	if found == false {
		depth++
		Bfs(newLinks, query, graph, final, depth)
	}
	fmt.Println(len(query))
	return graph
}

func Bfs2(links []string, query map[string]bool, graph *Graph, start string, final string) *Graph {

	shortDepth := 999999
	for len(links) > 0 {
		fmt.Println(len(query))
		currentLink := links[0]
		links = links[1:]

		links2, query2, graph2, found2 := getLinks(currentLink, query, graph, final)

		links = append(links, links2...)
		query = query2
		graph = graph2
		currentDepth := graph.maxDepth(start)
		fmt.Println(currentDepth)
		if currentDepth > shortDepth {
			break
		}
		if found2 == true {
			shortDepth = graph.maxDepth(start)
			fmt.Println(len(query2))
		}
	}
	return graph
}

func getLinks(html string, query map[string]bool, graph *Graph, final string) ([]string, map[string]bool, *Graph, bool) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(html)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	body, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	linkUrl := make(chan string)

	go func() {
		body.Find("p a[href]").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasPrefix(href, "/wiki/") && validLink(href) {
				link := "https://en.wikipedia.org" + href
				linkUrl <- link
			}
		})
		close(linkUrl)
	}()
	var links []string
	var linkFound []string
	for link := range linkUrl {
		if link == final {
			linkFound = append(linkFound, link)
			query[link] = true
			graph.AddEdge(html, link)
			return linkFound, query, graph, true
		}
		if query[link] == false {
			links = append(links, link)
			query[link] = true
			graph.AddEdge(html, link)
		}
	}

	// re := regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>`)
	// matches := re.FindAllStringSubmatch(string(body), -1)
	// uniqueLinks := make(map[string]bool)
	// for _, match := range matches {
	// 	if len(match) > 1 {
	// 		link := match[1]
	// 		if validLink(link) && !uniqueLinks[link] && !query[link] {
	// 			uniqueLinks[link] = true
	// 		}
	// 	}
	// }
	// var links []string
	// var linkFound []string
	// for link := range uniqueLinks {
	// 	link = "https://en.wikipedia.org" + link
	// 	if link == final {
	// 		linkFound = append(linkFound, link)
	// 		query[link] = true
	// 		graph.AddEdge(html, link)
	// 		return linkFound, query, graph, true
	// 	}

	// 	if query[link] == false {
	// 		links = append(links, link)
	// 		query[link] = true
	// 		graph.AddEdge(html, link)
	// 	}
	// }
	return links, query, graph, false
}

func validLink(link string) bool {
	if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, "(") && !strings.Contains(link, ".") && !strings.Contains(link, ",") && !strings.Contains(link, ":") && !strings.Contains(link, "#") && !strings.Contains(link, "%") && strings.ContainsAny(link, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return true
	}
	return false
}
