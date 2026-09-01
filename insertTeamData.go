package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InsertTeamsData(teamsData []TeamsD, season string) {
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

	cleanSeason := strings.TrimSuffix(season, ".csv")

	query := `
		INSERT INTO teams_season_data (
			competition_id, team_name, matches, wins, draws, loses,
			goals, goals_conceded, points, expected_goals,
			expected_goals_against, expected_points, season
		)
		VALUES (
			$1, $2, $3, $4,$5, $6, $7,$8, $9, $10, $11, $12, $13
		)
	`

	for _, team := range teamsData {
		_, err := conn.Exec(ctx, query,
			1,
			team.Team,
			team.Matches,
			team.Wins,
			team.Draws,
			team.Loses,
			team.Goals,
			team.GA,
			team.Points,
			team.XG,
			team.XGA,
			team.XPTS,
			cleanSeason,
		)

		if err != nil {
			// If one row fails, printing the team name helps you debug which row broke it
			fmt.Printf("Failed to insert team %s: %v\n", team.Team, err)
			panic(err) 
		}
	}


	fmt.Println("row inserted")
}