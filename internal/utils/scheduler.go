package utils

import ("time"
"database/sql"
"log"
"github.com/geshatude/JobBot/internal/utils"
)

// The Scheduler

func StartScheduler(db *sql.DB, interval time.Duration) {
	rnunCycle(db)
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        runCycle(db)
    }
}

func runCycle(db *sql.DB) {
    log.Println("Starting scrape cycle...")
    err := utils.Scrape(db)
    if err != nil {
        log.Println("Scrape error:", err)
        return
    }

    matches, err := utils.GetMatches(db)
    if err != nil {
        log.Println("Matching error:", err)
        return
    }

    err = utils.Notify(matches, db)
    if err != nil {
        log.Println("Notify error:", err)
        return
    }
    log.Println("Cycle complete.")
}