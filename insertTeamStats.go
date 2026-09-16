package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
)

func insertTeamStats(statsType string, season string, stats TeamStats) {
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

	// statsType := []string{"Attack speed", "Formation", "Game state", 
	// 	"Result", "Shot zones", "Situation", "Timing"}

	var tableName string

	switch statsType {
	case "Attack speed":
		tableName = "team_stats_by_attack_speed"
	case "Formation":
		tableName = "team_stats_by_formation"
	case "Game state":
		tableName = "team_stats_by_game_state"
	case "Result":
		tableName = ""
	case "Shot zones":
		tableName = "team_stats_by_shot_location"
	case "Situation":
		tableName = ""
	case "Timing":
		tableName = "team_stats_by_period"
	}
}