package main

import (
	"time"
	"encoding/csv"
	"os"
	"log"
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

func readCSV() ([]Match, map[string]bool) {
	files, err := os.ReadDir("./CSVs")
	if err != nil {
		log.Fatal(err)
	}

	fileMap := make(map[string]bool)

	for _, file := range files {
		if !file.IsDir() {
			fileMap[file.Name()] = false
		}
	}


}