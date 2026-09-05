package main

import (
		"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	// "strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	// "github.com/chromedp/cdproto/cdp"

)
// get all the teams stats form (https://understat.com/team/Sunderland/2025) 
// (8 files per team per season)

//next get players data from same site

//steps:  1 insert teams data in db 
// 2 create list of teams and scrap the teams data based on season
// 3 get player data (only 7 seasons)

func getTeamStats() {
	teams := getTeamsName()
	
	// Define the 7 tabs exactly as they appear on the screen
	// These must match the exact text inside the HTML <label> tags
	categories := []string{
		"Situation", 
		"Formation", 
		"Game state", 
		"Timing", 
		"Shot zones", 
		"Attack speed", 
		"Result",
	}

	wd, _ := os.Getwd()
	downloadDir := filepath.Join(wd, "teamsStats2018")
	os.MkdirAll(downloadDir, 0755)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc() 

	for _, team := range teams {
		func(currentTeam string) { 
			fmt.Printf("\n--- Processing team: %s ---\n", currentTeam)
			targetURL := "https://understat.com/team/" + currentTeam + "/2018" 

			ctx, cancelTab := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
			defer cancelTab()

			// Increased timeout because 7 downloads will take longer than 30 seconds
			ctx, cancelTimeout := context.WithTimeout(ctx, 2*time.Minute)
			defer cancelTimeout()

			// Set up channel with a small buffer
			done := make(chan string, 5)

			chromedp.ListenTarget(ctx, func(v interface{}) {
				if ev, ok := v.(*browser.EventDownloadProgress); ok {
					if ev.State == browser.DownloadProgressStateCompleted {
						done <- ev.GUID
					}
				}
			})

			// 1. Initial Setup and Navigate to the team's page
			err := chromedp.Run(ctx,
				browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
					WithDownloadPath(downloadDir).
					WithEventsEnabled(true),
				chromedp.Navigate(targetURL),
				// Wait for the download button to exist so we know the page loaded
				chromedp.WaitVisible(`button[name="table-export"]`, chromedp.ByQuery),
			)
			if err != nil {
				log.Printf("Failed to load page for %s: %v", currentTeam, err)
				return 
			}

			// 2. INNER LOOP: Iterate through the 7 categories
			for _, category := range categories {
				fmt.Printf("  -> Downloading %s...\n", category)

				var actions []chromedp.Action

				// ONLY click the tab if it's not the default "Situation" tab
				if category != "Situation" {
					
					// UPDATED: Look specifically for a <label> that has this exact text
					tabSelector := fmt.Sprintf(`//label[text()='%s']`, category)
					
					actions = append(actions, 
						chromedp.WaitVisible(tabSelector, chromedp.BySearch), // Good practice to wait for it first
						chromedp.Click(tabSelector, chromedp.BySearch),
						chromedp.Sleep(1 * time.Second), // Wait for table to load
					)
				}

				// ALWAYS add the download sequence
				actions = append(actions, 
					chromedp.WaitVisible(`button[name="table-export"]`, chromedp.ByQuery),
					chromedp.Click(`button[name="table-export"]`, chromedp.ByQuery),
					
					chromedp.WaitVisible(`button[data-ext="csv"]`, chromedp.ByQuery),
					chromedp.Click(`button[data-ext="csv"]`, chromedp.ByQuery),

					chromedp.WaitVisible(`i[class="popup-close fa fa-times"]`, chromedp.ByQuery),
					chromedp.Click(`i[class="popup-close fa fa-times"]`, chromedp.ByQuery),
					
					chromedp.Sleep(500 * time.Millisecond), 
				)

				// Execute the actions
				err = chromedp.Run(ctx, actions...)

				if err != nil {
					log.Printf("Failed to click through %s: %v", category, err)
					continue 
				}

				// 3. Wait for THIS specific download to finish
				select {
				case guid := <-done:
					time.Sleep(500 * time.Millisecond) // wait for file lock release

					oldFilePath := filepath.Join(downloadDir, guid)
					
					// Name the file dynamically: e.g., "Arsenal_Situation_2024.csv"
					newFileName := fmt.Sprintf("%s_%s_2018.csv", currentTeam, category)
					newFilePath := filepath.Join(downloadDir, newFileName) 

					if _, err := os.Stat(oldFilePath); err == nil {
						err := os.Rename(oldFilePath, newFilePath)
						if err != nil {
							log.Printf("Failed to rename file: %v", err)
						} else {
							fmt.Printf("     SUCCESS: Saved %s\n", newFileName)
						}
					}
				case <-time.After(15 * time.Second):
					// Use a local 15-second timeout per file, rather than context timeout
					log.Printf("     Timeout waiting for %s download\n", category)
				}
			}
		}(team)
	}
}



