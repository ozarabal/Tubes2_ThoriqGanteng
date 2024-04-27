package query

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Graph adalah tipe ada yang digunakan untuk menyimpan berbagai link
type Graph struct {
	adjacencyList map[string][]string
	visited       map[string]bool
}

// Fungsi untuk membuat graf baru
func NewGraph() *Graph {
	return &Graph{
		adjacencyList: make(map[string][]string),
		visited:       make(map[string]bool),
	}
}

// Fungsi untuk menambahkan sisi dalam graf dimana meminta parameter parent link dan child link
func (g *Graph) AddEdge(src, dest string) {
	g.adjacencyList[src] = append(g.adjacencyList[src], dest)
}

// Fungsi untuk mendapatkan kedalaman maksimal dalam graf
func (g *Graph) maxDepth(node string) int {
	g.visited = make(map[string]bool)
	return g.searchMax(node)
}

// Fungsi untuk mencari kedalaman maksimal graf
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

// Fungsi untuk mendapatkan semua path yang terhubung dari start hinga final
func GetAllPaths(graph *Graph, start string, final string, visited map[string]bool, path []string, allPaths *[][]string) {
	visited[start] = true
	path = append(path, start)
	if start == final {
		*allPaths = append(*allPaths, append([]string{}, path...))
	} else {
		for _, neighbor := range graph.adjacencyList[start] {
			if !visited[neighbor] {
				GetAllPaths(graph, neighbor, final, visited, path, allPaths)
			}
		}
	}
	path = path[:len(path)-1]
	visited[start] = false
}

// Fungsi untuk mendapatkan path terpendek dari start hingga final dengan menggunakan algoritma BFS
func Bfs2(queueLinks []string, visitedLink map[string]bool, graph *Graph, start string, final string) *Graph {
	// Menginisiasi kedalaman terpendek dan memulai waktu pencarian
	shortDepth := 999999
	timeoutSeconds := 290
	timeout := time.Duration(timeoutSeconds) * time.Second
	startTime := time.Now()

	// Melakukan pencarian berdasarkan dengan queueLinks
	for len(queueLinks) > 0 {

		// Jika waktu sudah melewati batas waktu tertentu maka akan mengembalikan hasil yang sudah ada
		if time.Since(startTime) > timeout {
			return graph
		}

		// Melakukan pengambilan seluruh hyperlink yang ada dalam suatu page
		currentLink := queueLinks[0]
		queueLinks = queueLinks[1:]
		links2, query2, graph2, found2 := getLinks(currentLink, visitedLink, graph, final)
		queueLinks = append(queueLinks, links2...)
		visitedLink = query2
		graph = graph2

		// Melakukan pengecekan apakah kedalaman sekarang sudah melewati kedalaman terpendak
		currentDepth := graph.maxDepth(start)
		if currentDepth > shortDepth {
			break
		}

		// Jika link final ketemu, maka kedalaman sekarang akan dijadikan sebagai kedalaman terpendek
		if found2 == true {
			shortDepth = graph.maxDepth(start)
		}
	}
	return graph
}

// Fungsi untuk mendapatkan berbagai hyperlink yang ada di suatu page
func getLinks(html string, visitedLink map[string]bool, graph *Graph, final string) ([]string, map[string]bool, *Graph, bool) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(html)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	re := regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	uniqueLinks := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]
			if validLink(link) && !uniqueLinks[link] && !visitedLink[link] {
				uniqueLinks[link] = true
			}
		}
	}
	var newLinks []string
	var linkFound []string
	for link := range uniqueLinks {
		link = "https://en.wikipedia.org" + link
		if link == final {
			linkFound = append(linkFound, link)
			visitedLink[link] = true
			graph.AddEdge(html, link)
			return linkFound, visitedLink, graph, true
		}

		if visitedLink[link] == false {
			newLinks = append(newLinks, link)
			visitedLink[link] = true
			graph.AddEdge(html, link)
		}
	}
	return newLinks, visitedLink, graph, false
}

// Fungsi untuk mengecek apakah link tersebut termasuk link yang valid atau tidak
func validLink(link string) bool {
	matched, _ := regexp.MatchString(`^/wiki/[A-Z][^(),:%#]*$`, link)
	return matched && !strings.Contains(link, ".")
}
