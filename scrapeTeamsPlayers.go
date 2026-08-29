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

func GetTeamAndPlayerData(targetURL string, season string) {
	
	// 2. Set up a temporary directory for the download
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	downloadDir := filepath.Join(wd, "downloads")
	os.MkdirAll(downloadDir, 0755)

	// 3. Create context and allocate the browser
	// Note: You can add chromedp.Flag("headless", false) to watch it happen visually for debugging
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), 
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Add a timeout so the script doesn't hang forever if something fails
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 4. Set up a channel to listen for the download completion
	done := make(chan string, 1)

	// Listen for browser events
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			// When the download state becomes "completed", send the GUID to our channel
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
			}
		}
	})

	fmt.Println("Navigating to page and clicking buttons...")

	// 5. Execute the browser actions
	err = chromedp.Run(ctx,
		// Tell the browser to allow downloads and where to put them
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(downloadDir).
			WithEventsEnabled(true),

		// Navigate to the target page
		chromedp.Navigate(targetURL),

		// Wait for the FIRST export button to be visible, then click it
		// Using the selector provided: button[name="table-export"]
		chromedp.WaitVisible(`button[name="table-export"]`, chromedp.ByQuery),
		chromedp.Click(`button[name="table-export"]`, chromedp.ByQuery),

		// Wait for the popup and the CSV button to be visible, then click it
		// Using the selector provided: button[data-ext="csv"]
		chromedp.WaitVisible(`button[data-ext="csv"]`, chromedp.ByQuery),
		chromedp.Click(`button[data-ext="csv"]`, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("Failed to execute chromedp actions: %v", err)
	}

	// 6. Wait for the download to finish
	fmt.Println("Waiting for download to complete...")
	
	select {
	case guid := <-done:
		fmt.Printf("Download completed. GUID: %s\n", guid)

		// Give the OS a tiny fraction of a second to release the file lock
		time.Sleep(500 * time.Millisecond)

		// 1. The old file is exactly the GUID name
		oldFilePath := filepath.Join(downloadDir, guid)
		
		// 2. The new file is your season year. 
		// IMPORTANT: Add +".csv" here so your computer knows it is a spreadsheet!
		newFilePath := filepath.Join(downloadDir, season+".csv") 

		// 3. Check if the GUID file exists, then rename it
		if _, err := os.Stat(oldFilePath); err == nil {
			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				log.Fatalf("Failed to rename file: %v", err)
			}
			fmt.Printf("SUCCESS: File renamed to: %s\n", newFilePath)
		} else {
			fmt.Printf("Error: Expected to find file %s but it wasn't there.\n", oldFilePath)
		}

	case <-ctx.Done():
		log.Fatal("Timeout waiting for download to complete")
	}
}


func GetPlayerData(targetURL string, season string) {

	// 2. Set up a temporary directory for the download
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	downloadDir := filepath.Join(wd, "downloadPlayers")
	os.MkdirAll(downloadDir, 0755)

	// 3. Create context and allocate the browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Add a timeout so the script doesn't hang forever if something fails
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 4. Set up a channel to listen for the download completion
	done := make(chan string, 1)

	// Listen for browser events
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
			}
		}
	})

	fmt.Println("Navigating to page and clicking buttons...")

	// 5. Execute the browser actions
	err = chromedp.Run(ctx,
		// Tell the browser to allow downloads and where to put them
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(downloadDir).
			WithEventsEnabled(true),

		// Navigate to the target page
		chromedp.Navigate(targetURL),

		// --- CHANGED SECTION ---
		// Wait for the buttons to appear on the page
		chromedp.WaitVisible(`button[name="table-export"]`, chromedp.ByQuery),
		
		// Use JavaScript to select all matching buttons and click the second one (index [1])
		chromedp.EvaluateAsDevTools(`document.querySelectorAll('button[name="table-export"]')[1].click();`, nil),
		// -----------------------

		// Wait for the popup and the CSV button to be visible, then click it
		chromedp.WaitVisible(`button[data-ext="csv"]`, chromedp.ByQuery),
		chromedp.Click(`button[data-ext="csv"]`, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("Failed to execute chromedp actions: %v", err)
	}

	// 6. Wait for the download to finish
	fmt.Println("Waiting for download to complete...")

	select {
	case guid := <-done:
		fmt.Printf("Download completed. GUID: %s\n", guid)

		// Give the OS a tiny fraction of a second to release the file lock
		time.Sleep(500 * time.Millisecond)

		// 1. The old file is exactly the GUID name
		oldFilePath := filepath.Join(downloadDir, guid)
		
		// 2. The new file is your season year. 
		// IMPORTANT: Add +".csv" here so your computer knows it is a spreadsheet!
		newFilePath := filepath.Join(downloadDir, season+".csv") 

		// 3. Check if the GUID file exists, then rename it
		if _, err := os.Stat(oldFilePath); err == nil {
			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				log.Fatalf("Failed to rename file: %v", err)
			}
			fmt.Printf("SUCCESS: File renamed to: %s\n", newFilePath)
		} else {
			fmt.Printf("Error: Expected to find file %s but it wasn't there.\n", oldFilePath)
		}

	case <-ctx.Done():
		log.Fatal("Timeout waiting for download to complete")
	}
}