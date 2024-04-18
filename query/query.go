package query

import (
    "io"
	"net/http"
    "fmt"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

// Node adalah tipe data untuk node dalam struktur data graph
type Node struct {
	PageURL  string   // URL dari halaman Wikipedia
	Children []*Node // Anak-anak dari node
	Parent   *Node    // Parent dari node
}

func FindTree(start, goal string) ([]*Node, error) {
	startURL := "https://en.wikipedia.org/wiki/" + start
	goalURL := "https://en.wikipedia.org/wiki/" + goal
	fmt.Println("Start URL:", startURL)
	paths, err := ShortestPath(startURL, goalURL)
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

// ShortestPath mencari jalur terpendek dari startURL ke goalURL
func ShortestPath(startURL, goalURL string) ([][]*Node, error) {
	root := &Node{PageURL: startURL} // Membuat node root dengan startURL sebagai PageURL
	visited := make(map[string]bool) // Membuat map untuk menyimpan node yang sudah dieksplorasi

	paths := [][]*Node{} // Menyimpan semua jalur yang terbentuk
	paths = getPaths(startURL, goalURL, visited, root, paths)

	if len(paths) == 0 {
		return nil, fmt.Errorf("shortest path not found")
	}

	return paths, nil
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

// getPaths mengambil jalur-jalur yang dapat dicapai dari startURL ke goalURL
func getPaths(startURL, goalURL string, visited map[string]bool, parent *Node, paths [][]*Node) [][]*Node {
	if parent.PageURL == goalURL { // Jika node saat ini adalah goal node
		// Path ditemukan, tambahkan jalur ke paths
		paths = append(paths, getPath(parent))
		return paths
	}

	if visited[parent.PageURL] { // Jika node sudah dieksplorasi sebelumnya
		return paths
	}
	visited[parent.PageURL] = true // Tandai node sebagai sudah dieksplorasi

	// Ambil tautan-tautan dari halaman node
	links, err := GetLinks(parent.PageURL)
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return paths
	}

	// Untuk setiap tautan, buat node anak dan rekursif cari jalur
	for _, link := range links {
		child := &Node{PageURL: link, Parent: parent}
		parent.Children = append(parent.Children, child)
		paths = getPaths(startURL, goalURL, visited, child, paths)
	}

	return paths
}


// getLinks mengambil tautan-tautan dari halaman Wikipedia
func GetLinks(pageURL string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0") // Set header User-Agent untuk menghindari pemblokiran

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
		}
	})

	return links, nil
}