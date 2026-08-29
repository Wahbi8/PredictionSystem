package main

import (
	"strconv"
)

func main() {
	// scrape data
	// GetCSVFiles()

	// insert the scraped data
	// readCSV()


	// scrape players and teams data
	for year := 2014; year <= 2025; year++ {
		url := "https://understat.com/league/EPL/" + strconv.Itoa(year)
		GetTeamAndPlayerData(url, strconv.Itoa(year))
		GetPlayerData(url, strconv.Itoa(year))

	}

	// for year := 2014; year <= 2025; year++ {
	// 	url := "https://understat.com/league/EPL/" + strconv.Itoa(year)
	// 	GetPlayerData(url)
	// }
}