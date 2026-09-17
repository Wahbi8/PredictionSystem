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

type TeamStats struct {
	ID            int64
	TeamName      string
	CompetitionID int
	Season        string
	Statistic     string
	Min           int
	Shots         int
	Goals         int
	ShotsAgainst  int
	GoalsAgainst  int
	XG            float64
	XGA           float64
	XGDiff        float64
	XGPerSh       float64
	XGAPerSh      float64
	XG90          float64
	XGA90         float64
}

func readTeamsStats() {
	files, err := os.ReadDir("TeamsStats2018")
	if err != nil {
		log.Fatal(err)
	}

	// statsType := []string{"Attack speed", "Formation", "Game state", 
	// 	"Result", "Shot zones", "Situation", "Timing"}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		path := filepath.Join("./TeamsStats2018", file.Name())

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		reader.Comma = ';'
		reader.FieldsPerRecord = -1
		reader.LazyQuotes = true

		// Skip header
		_, err = reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		s := strings.Split(file.Name(),"_")
		statType := s[1]

		stats := []TeamStats{}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if len(record) < 11 {
				continue
			}

			stat := TeamStats{
				ID:        parseID(record[0]),
				Statistic: record[1],
				TeamName: s[0],
			}

			offset := 0

			switch statType {
			case "Formation", "Game state":
				offset = 1
				stat.Min = parseInt(record[2])
				stat.XG90 = parseFloat(record[10])
				stat.XGA90 = parseFloat(record[11])
				

			case "Attack speed", "Result", "Shot zones", "Situation", "Timing":
				offset = 0
				stat.XGPerSh = parseFloat(record[9])
				stat.XGAPerSh = parseFloat(record[10])
			}

			// Shared fields parsed once using the offset
			stat.Shots = parseInt(record[2+offset])
			stat.Goals = parseInt(record[3+offset])
			stat.ShotsAgainst = parseInt(record[4+offset])
			stat.GoalsAgainst = parseInt(record[5+offset])
			stat.XG = parseFloat(record[6+offset])
			stat.XGA = parseFloat(record[7+offset])
			stat.XGDiff = parseFloat(record[8+offset])

			stats = append(stats, stat)
		}

		f.Close() 

		
		// db insert func to be added
		insertTeamStats(s[1], s[2], stats)

		newPath := filepath.Join("processedTeamStatsafter25", file.Name())
		err = os.Rename(path, newPath)
		if err != nil {
			log.Fatalf("Failed to move file: %v", err)
		}
		
		fmt.Printf("Processed and moved file: %v\n", file.Name())

		fmt.Printf("%s: %d stats\n", file.Name(), len(stats))
	}
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseID(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}