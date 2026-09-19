package main

import ("database/sql"
"log"
"github.com/geshatude/JobBot/utils"
_"github.com/lib/pq")

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:mysecretpassword@127.0.0.1:5432/jobbot?sslmode=disable")
	if err != nil {
		log.Fatal("database connection error", err)
	}
	err = utils.Scrape(db)
	if err != nil {
		log.Fatal("scraping error", err)
	}
}