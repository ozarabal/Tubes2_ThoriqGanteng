package query

import (
	"fmt"
	// "io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	// "sync"
	"time"
	// "github.com/PuerkitoBio/goquery"
)

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

func Bfs2(links []string, query map[string]bool, graph *Graph, start string, final string) *Graph {
	shortDepth := 999999
	timeoutSeconds := 290
	timeout := time.Duration(timeoutSeconds) * time.Second
	startTime := time.Now()
	for len(links) > 0 {
		if time.Since(startTime) > timeout {
			return graph
		}
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
	body, err := ioutil.ReadAll(resp.Body)
	//body, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// linkUrl := make(chan string)

	// go func() {
	// 	body.Find("p a[href]").Each(func(i int, s *goquery.Selection) {
	// 		href, exists := s.Attr("href")
	// 		if exists && strings.HasPrefix(href, "/wiki/") && validLink(href) {
	// 			link := "https://en.wikipedia.org" + href
	// 			linkUrl <- link
	// 		}
	// 	})
	// 	close(linkUrl)
	// }()
	// var links []string
	// var linkFound []string
	// for link := range linkUrl {
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

	re := regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	uniqueLinks := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]
			if validLink(link) && !uniqueLinks[link] && !query[link] {
				uniqueLinks[link] = true
			}
		}
	}
	var links []string
	var linkFound []string
	for link := range uniqueLinks {
		link = "https://en.wikipedia.org" + link
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
	return links, query, graph, false
}

func validLink(link string) bool {
	if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, "(") && !strings.Contains(link, ".") && !strings.Contains(link, ",") && !strings.Contains(link, ":") && !strings.Contains(link, "#") && !strings.Contains(link, "%") && strings.ContainsAny(link, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return true
	}
	return false
}
