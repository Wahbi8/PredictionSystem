package main

import (
	"os"
	"net/http"
	"log"
	"io"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func GetTeamAndPlayerData() {
	url := "https://understat.com/league/EPL/2014"

	fmt.Println("connecting to website...")

	req, err := http.NewRequest("GEt", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("start scraping link")

	doc.Find()

}