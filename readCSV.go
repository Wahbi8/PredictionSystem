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

// dates red (2002, 2003, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026)

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
	if s == "" {
		return time.Time{}
	}

	chunks := strings.Split(s, "/")
	d := s

	if len(chunks) == 3 {
		switch chunks[2] {
		case "00":
			d = chunks[0] + "/" + chunks[1] + "/2000"
		case "01":
			d = chunks[0] + "/" + chunks[1] + "/2001"
		case "02":
			d = chunks[0] + "/" + chunks[1] + "/2002"
		case "03":
			d = chunks[0] + "/" + chunks[1] + "/2003"
		case "04":
			d = chunks[0] + "/" + chunks[1] + "/2004"
		case "05":
			d = chunks[0] + "/" + chunks[1] + "/2005"
		case "06":
			d = chunks[0] + "/" + chunks[1] + "/2006"
		case "07":
			d = chunks[0] + "/" + chunks[1] + "/2007"
		case "08":
			d = chunks[0] + "/" + chunks[1] + "/2008"
		case "09":
			d = chunks[0] + "/" + chunks[1] + "/2009"
		case "10":
			d = chunks[0] + "/" + chunks[1] + "/2010"
		case "11":
			d = chunks[0] + "/" + chunks[1] + "/2011"
		case "12":
			d = chunks[0] + "/" + chunks[1] + "/2012"
		case "13":
			d = chunks[0] + "/" + chunks[1] + "/2013"
		case "14":
			d = chunks[0] + "/" + chunks[1] + "/2014"
		case "15":
			d = chunks[0] + "/" + chunks[1] + "/2015"
		case "16":
			d = chunks[0] + "/" + chunks[1] + "/2016"
		case "17":
			d = chunks[0] + "/" + chunks[1] + "/2017"
		case "18":
			d = chunks[0] + "/" + chunks[1] + "/2018"
		case "19":
			d = chunks[0] + "/" + chunks[1] + "/2019"
		case "20":
			d = chunks[0] + "/" + chunks[1] + "/2020"
		case "21":
			d = chunks[0] + "/" + chunks[1] + "/2021"
		case "22":
			d = chunks[0] + "/" + chunks[1] + "/2022"
		case "23":
			d = chunks[0] + "/" + chunks[1] + "/2023"
		case "24":
			d = chunks[0] + "/" + chunks[1] + "/2024"
		case "25":
			d = chunks[0] + "/" + chunks[1] + "/2025"
		case "26":
			d = chunks[0] + "/" + chunks[1] + "/2026"
		}
	}

	t, _ := time.Parse("02/01/2006", d)
	return t
}

