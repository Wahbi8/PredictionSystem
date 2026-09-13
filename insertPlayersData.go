package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func insertPlayersData(season string, players []Player) {
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
		INSERT INTO player_stats (
			player,
			team,
			apps,
			min,
			goals,
			a,
			xG,
			xA,
			xG90,
			xA90,
			season,
			competition_id
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12
		)
	`

	for _, player := range players {
		_, err := conn.Exec(ctx, query,
			player.Name,
			player.Team,
			player.Apps,
			player.Min,
			player.Goals,
			player.A,
			player.XG,
			player.XA,
			player.XG90,
			player.XA90,
			cleanSeason,
			1,
		)

		if err != nil {
			// If one row fails, printing the team name helps you debug which row broke it
			fmt.Printf("Failed to insert team %s: %v\n", &player.Name, err)
			panic(err) 
		}
	}


	fmt.Println("row inserted")
}