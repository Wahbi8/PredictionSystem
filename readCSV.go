package main

import (
	"encoding/csv"
	"log"
	"os"
	"io"
	"path/filepath"
	"time"
	"strconv"
	"fmt"
	"strings"
)

type Match struct {
	fileName string
	Div      string    `csv:"Div"`
	Date     time.Time `csv:"Date"`
	HomeTeam string    `csv:"HomeTeam"`
	AwayTeam string    `csv:"AwayTeam"`

	FTHG int    `csv:"FTHG"`
	FTAG int    `csv:"FTAG"`
	FTR  string `csv:"FTR"`

	HTHG int    `csv:"HTHG"`
	HTAG int    `csv:"HTAG"`
	HTR  string `csv:"HTR"`

	Referee string `csv:"Referee"`

	HS  int `csv:"HS"`
	AS  int `csv:"AS"`
	HST int `csv:"HST"`
	AST int `csv:"AST"`
	HF  int `csv:"HF"`
	AF  int `csv:"AF"`
	HC  int `csv:"HC"`
	AC  int `csv:"AC"`
	HY  int `csv:"HY"`
	AY  int `csv:"AY"`
	HR  int `csv:"HR"`
	AR  int `csv:"AR"`

	B365H float64 `csv:"B365H"`
	B365D float64 `csv:"B365D"`
	B365A float64 `csv:"B365A"`
}

func readCSV() {
	files, err := os.ReadDir("./CSVs")
	if err != nil {
		log.Fatal(err)
	}


	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		path := filepath.Join("./CSVs", file.Name())

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		reader := csv.NewReader(f)
		reader.FieldsPerRecord = -1 
		_, err = reader.Read() // skip header
		if err != nil {
			log.Fatal(err)
		}

		matches := []Match{}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if len(record) < 26 {
				continue
			}

			for i := range record {
				// Replace Windows non-breaking space (0xa0) with a normal space
				record[i] = strings.ReplaceAll(record[i], "\xa0", " ")
				// Strip out any other invalid UTF-8 characters just to be safe
				record[i] = strings.ToValidUTF8(record[i], "")
			}

			matches = append(matches, Match{
				fileName: file.Name(),
				Div:      record[0],
				Date:     parseDate(record[1]),
				HomeTeam: record[2],
				AwayTeam: record[3],

				FTHG: atoi(record[4]),
				FTAG: atoi(record[5]),
				FTR:  record[6],

				HTHG: atoi(record[7]),
				HTAG: atoi(record[8]),
				HTR:  record[9],

				Referee: record[10],

				HS: atoi(record[11]),
				AS: atoi(record[12]),
				HST: atoi(record[13]),
				AST: atoi(record[14]),
				HF: atoi(record[15]),
				AF: atoi(record[16]),
				HC: atoi(record[17]),
				AC: atoi(record[18]),
				HY: atoi(record[19]),
				AY: atoi(record[20]),
				HR: atoi(record[21]),
				AR: atoi(record[22]),

				B365H: atof(record[23]),
				B365D: atof(record[24]),
				B365A: atof(record[25]),
			})			
		}
		f.Close() 

		insertCSVToDb(matches)

		// MOVE THE FILE 
		newPath := filepath.Join("processedCSVs", file.Name())
		err = os.Rename(path, newPath)
		if err != nil {
			log.Fatalf("Failed to move file: %v", err)
		}
		
		fmt.Printf("Processed and moved file: %v\n", file.Name())
	}

	
}

func atoi(s string) int {
    v, _ := strconv.Atoi(s)
    return v
}

func atof(s string) float64 {
    v, _ := strconv.ParseFloat(s, 64)
    return v
}

func parseDate(s string) time.Time {
    t, _ := time.Parse("02/01/2006", s)
    return t
}