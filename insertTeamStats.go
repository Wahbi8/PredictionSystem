package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

)

func insertTeamStats(statsType string, season string, stats []TeamStats) {
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
		tableName = "team_stats_by_shot_outcome"
	case "Shot zones":
		tableName = "team_stats_by_shot_location"
	case "Situation":
		tableName = "team_stats_by_chance_type"
	case "Timing":
		tableName = "team_stats_by_period"
	default:
		return
	}

	var query string

	if statsType == "Formation" || statsType == "Game state" {
		query = fmt.Sprintf(`
			INSERT INTO %s (
				team_name, competition_id, season, statistic, min,
				shots, goals, shots_against, goals_against,
				xg, xga, xg_diff, xg90, xga90
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, tableName)

		for _, stat := range stats {
			_, err = conn.Exec(ctx, query,
				stat.TeamName,
				stat.CompetitionID,
				cleanSeason,
				stat.Statistic,
				stat.Min,
				stat.Shots,
				stat.Goals,
				stat.ShotsAgainst,
				stat.GoalsAgainst,
				stat.XG,
				stat.XGA,
				stat.XGDiff,
				stat.XG90,
				stat.XGA90,
			)
		}
	} else {
		query = fmt.Sprintf(`
			INSERT INTO %s (
				team_name, competition_id, season, statistic,
				shots, goals, shots_against, goals_against,
				xg, xga, xg_diff, xgpersh, xgapersh
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, tableName)

		for _, stat := range stats {
			_, err = conn.Exec(ctx, query,
				stat.TeamName,
				stat.CompetitionID,
				cleanSeason,
				stat.Statistic,
				stat.Shots,
				stat.Goals,
				stat.ShotsAgainst,
				stat.GoalsAgainst,
				stat.XG,
				stat.XGA,
				stat.XGDiff,
				stat.XGPerSh,
				stat.XGAPerSh,
			)
		}
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("data inserted")
}