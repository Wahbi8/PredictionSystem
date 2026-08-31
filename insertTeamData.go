package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InsertTeamsData(teamsData []TeamsD) {
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

	query = `
		INSERT INTO teams_season_data (
			competition_id, team_name, matches, wins, draws, loses,
			goals, goals_conceded, points, expected_goalsm,
			expected_goals_against, expected_points 
		)
		VALUES (
			$1, $2, $3, $4,$5, $6, $7,$8, $9, $10, $11, $12
		)
	`

	for _, team := teamsData {
		_, err := conn.Exec(ctx, query,
		1,
		
		)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("row inserted")
}