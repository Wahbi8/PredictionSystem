package main

import (
	"os"
	"io"
	"encoding/csv"
	"path/filepath"
	"log"
	"strings"
	"strconv"
	"fmt"
)

type TeamsD struct{
	Team    string  `json:"team"`
	Matches int     `json:"matches"`
	Wins    int     `json:"wins"`
	Draws   int     `json:"draws"`
	Loses   int     `json:"loses"`
	Goals   int     `json:"goals"`
	GA      int     `json:"ga"`
	Points  int     `json:"points"`
	XG      float64 `json:"xG"`
	XGA     float64 `json:"xGA"`
	XPTS    float64 `json:"xPTS"`
	season 	string
}

func readTeamsFromCsv() {
	files, err := os.ReadDir("./downloads")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println("looking for files")
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}
		
		path := filepath.Join("./downloads", file.Name())

		f, err := os.Open(path) 
		if err != nil {
			log.Fatal(err)
		}

		reader := csv.NewReader(f)
		
		// ---> THE MAGIC FIX IS HERE <---
		reader.Comma = ';' 
		
		reader.FieldsPerRecord = -1
		reader.LazyQuotes = true

		// Skip header
		_, err = reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		teamd := []TeamsD{}

		// --- ROW READING LOOP ---
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break // End of file reached
			}
			if err != nil {
				log.Fatal(err)
			}

			if len(record) < 12 {
				continue
			}

			for i := range record {
				record[i] = strings.ReplaceAll(record[i], "\xa0", " ")
				record[i] = strings.ToValidUTF8(record[i], "")
			}

			matches, _ := strconv.Atoi(record[2])
			wins, _ := strconv.Atoi(record[3])
			draws, _ := strconv.Atoi(record[4])
			loses, _ := strconv.Atoi(record[5])
			goals, _ := strconv.Atoi(record[6])
			ga, _ := strconv.Atoi(record[7])
			points, _ := strconv.Atoi(record[8])
			xg, _ := strconv.ParseFloat(record[9], 64)
			xga, _ := strconv.ParseFloat(record[10], 64)
			xpts, _ := strconv.ParseFloat(record[11], 64)

			teamd = append(teamd, TeamsD{
				Team:    record[1],
				Matches: matches,
				Wins:    wins,
				Draws:   draws,
				Loses:   loses,
				Goals:   goals,
				GA:      ga,
				Points:  points,
				XG:      xg,
				XGA:     xga,
				XPTS:    xpts,
				season:  strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			})
		}
		// --- END OF ROW READING LOOP ---

		f.Close() 

		// Now teamd will actually have 20 teams in it!
		InsertTeamsData(teamd, file.Name())

		newPath := filepath.Join("processedTeamCSV", file.Name())
		err = os.Rename(path, newPath)
		if err != nil {
			log.Fatalf("Failed to move file: %v", err)
		}
		
		fmt.Printf("Processed and moved file: %v\n", file.Name())
	}
}