package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var urls = []string{
	"https://go.dev",
	"https://golang.org",
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("===== EXERCÍCIO 1: Título + Meta Description =====")
	for _, u := range urls {
		doc, err := fetchDocument(client, u)
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		title := doc.Find("title").Text()
		description, _ := doc.Find(`meta[name="description"]`).Attr("content")
		if description == "" {
			// tenta meta description genérica (alguns sites usam og:description)
			description, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
		}

		fmt.Println("🔍", u)
		fmt.Println("  📄 Title:      ", title)
		fmt.Println("  📝 Description:", description)
		fmt.Println()
	}

	fmt.Println("===== EXERCÍCIO 2: Contagem de Links =====")
	for _, u := range urls {
		doc, err := fetchDocument(client, u)
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		linksCount := doc.Find("a").Length()
		fmt.Printf("🔗 %s -> %d links\n", u, linksCount)
	}

	fmt.Println("\n===== EXERCÍCIO 3: Timeout diferente por URL =====")
	for i, u := range urls {
		timeout := time.Duration(2+i) * time.Second
		c := &http.Client{Timeout: timeout}
		fmt.Printf("⏳ Timeout %s - Acessando: %s\n", timeout, u)
		_, err := fetchDocument(c, u)
		if err != nil {
			log.Println("Erro (timeout/url):", err)
		} else {
			fmt.Println("✅ Requisição concluída com sucesso!")
		}
	}
}

// fetchDocument faz a requisição HTTP e retorna um *goquery.Document.
// Fechamos resp.Body com defer dentro da função (seguro e limpo).
func fetchDocument(client *http.Client, url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("criar request: %w", err)
	}
	// User-Agent básico para evitar bloqueios simples
	req.Header.Set("User-Agent", "Golang Web Crawler - aula (github.com/PuerkitoBio/goquery)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro client.Do: %w", err)
	}
	defer resp.Body.Close()

	// checar status http
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status code não OK: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery parse: %w", err)
	}

	return doc, nil
}
