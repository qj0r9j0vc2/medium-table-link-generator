package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	urlFlag := flag.String("url", "", "Medium article URL")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go -url <medium-article-url>")
		os.Exit(1)
	}

	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	parsed, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}
	baseURL := fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)

	doc.Find("h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		text := strings.TrimSpace(s.Text())
		link := fmt.Sprintf("%s#%s", baseURL, id)
		fmt.Printf("%s %s\n", text, link)
	})
}
