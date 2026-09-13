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



func readTeamsStats() {
	files, err := os.ReadDir("teamsStats")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		path := filepath.Join("./teamsStats", file.Name())

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

		players := []Player{}

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

			for i := range record {
				record[i] = strings.ReplaceAll(record[i], "\xa0", " ")
				record[i] = strings.ToValidUTF8(record[i], "")
				record[i] = strings.TrimSpace(record[i])
			}

			apps, _ := strconv.Atoi(record[3])
			min, _ := strconv.Atoi(record[4])
			goals, _ := strconv.Atoi(record[5])
			assists, _ := strconv.Atoi(record[6])

			xG, _ := strconv.ParseFloat(record[7], 64)
			xA, _ := strconv.ParseFloat(record[8], 64)
			xG90, _ := strconv.ParseFloat(record[9], 64)
			xA90, _ := strconv.ParseFloat(record[10], 64)

			player := Player{
				Name: record[1],
				Team:   record[2],
				Apps:   apps,
				Min:    min,
				Goals:  goals,
				A:      assists,
				XG:     xG,
				XA:     xA,
				XG90:   xG90,
				XA90:   xA90,
			}

			players = append(players, player)
		}

		f.Close() 

		insertPlayersData(file.Name(), players)

		newPath := filepath.Join("processedTeamStats", file.Name())
		err = os.Rename(path, newPath)
		if err != nil {
			log.Fatalf("Failed to move file: %v", err)
		}
		
		fmt.Printf("Processed and moved file: %v\n", file.Name())
	}
}