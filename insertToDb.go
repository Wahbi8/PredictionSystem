package main

import (
	"os"
	"github.com/joho/godotenv"
	// "fmt"
)

func insertCSVToDb() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	connectionStr := os.Getenv("connection")

	// read csv files 
	
}