package main

import (
	"os"
	"context"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)

func getTeamsName() []string {
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

	query := `SELECT DISTINCT team_name
		FROM teams_season_data
		WHERE season = '2018'`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var teams []string

	for rows.Next() {
		var teamName string

		err := rows.Scan(&teamName)
		if err != nil {
			panic(err)
		}

		teams = append(teams, teamName)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return teams
}