package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	// "path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)



func GetCSVFiles() {
	url := "https://www.football-data.co.uk/englandm.php"

	fmt.Println("Connecting to website...")

	// 1. Create a custom request so we can add a User-Agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set the User-Agent to mimic a real web browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 3. Check the status code! Ensure we actually got the webpage.
	if resp.StatusCode != 200 {
		log.Fatalf("Error: Status code %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start scraping links...")
	
	count := 0
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		linkText := strings.TrimSpace(s.Text())

		if linkText == "Premier League" {
			href, exists := s.Attr("href")
			if !exists || len(href) == 0 {
				return
			}

			fmt.Printf("Found Link: %s\n", href)
			DownloadFile(href)
			count++
		}
	})
	
	fmt.Printf("\nFinished! Downloaded %d files.\n", count)
}

func DownloadFile(href string) {
	url := "https://www.football-data.co.uk/" + href

	name := strings.ReplaceAll(href, "/", "_")
	// fileName := filepath.Base(href)

	fmt.Printf(" -> Downloading: %s\n", name)

	// Create request with User-Agent for downloading the file as well
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v\n", name, err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to download %s: %v\n", name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Failed to download %s. Status Code: %d\n", name, resp.StatusCode)
		return
	}

	file, err := os.Create("CSVs/" + name)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}
}
