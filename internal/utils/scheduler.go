package utils

import ("time"
"database/sql"
"log"
)

// The Scheduler

func StartScheduler(db *sql.DB, interval time.Duration) {
	runCycle(db)
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        runCycle(db)
    }
}

func runCycle(db *sql.DB) {
    log.Println("Starting scrape cycle...")
    err := Scrape(db)
    if err != nil {
        log.Println("Scrape error:", err)
        return
    }
    log.Println("Scrape cycle completed. Starting matching and notification...")
    matches, err := GetMatches(db)
    if err != nil {
        log.Println("Matching error:", err)
        return
    }
    log.Printf("Found %d matches. Sending notifications...", len(matches))
    err = Notify(matches, db)
    if err != nil {
        log.Println("Notify error:", err)
        return
    }
    log.Println("Cycle complete.")
}