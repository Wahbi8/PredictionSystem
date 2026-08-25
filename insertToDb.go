package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	// "log"
	// "io"
	// "path/filepath"

	"github.com/jackc/pgx/v5"
)

func insertCSVToDb(matches []Match) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	connectionStr := os.Getenv("connection")

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connectionStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	query := `
		INSERT INTO matches (
			competition_id, "Date", "HomeTeam", "AwayTeam",
			"FTHG", "FTAG", "FTR",
            "HTHG", "HTAG", "HTR",
            "Referee",
            "HS", "AS", "HST", "AST", "HF", "AF", "HC", "AC", "HY", "AY", "HR", "AR",
            "B365H", "B365D", "B365A"
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10,
			$11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23,
			$24, $25, $26
		)
`

	for _, match := range matches {
		_, err := conn.Exec(ctx, query,
			1,
			match.Date,
			match.HomeTeam,
			match.AwayTeam,
			match.FTHG,
			match.FTAG,
			match.FTR,
			match.HTHG,
			match.HTAG,
			match.HTR,
			match.Referee,
			match.HS,
			match.AS,
			match.HST,
			match.AST,
			match.HF,
			match.AF,
			match.HC,
			match.AC,
			match.HY,
			match.AY,
			match.HR,
			match.AR,
			match.B365H,
			match.B365D,
			match.B365A,
		)

		if err != nil {
			panic(err)
		}

		fmt.Print("row inserted \n")
	}
	
}