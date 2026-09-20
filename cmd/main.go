package main

import ("database/sql"
"log"
"github.com/geshatude/JobBot/internal/server"
_"github.com/lib/pq"
"github.com/geshatude/JobBot/internal/utils"
"net/http"
"time"
"log"
"os")

func main() {
	mymux := http.NewServeMux()
	db, err := sql.Open("postgres", "postgres://postgres:mysecretpassword@127.0.0.1:5432/jobbot?sslmode=disable")
	if err != nil {
		log.Fatal("database connection error", err)
	}
	go utils.StartScheduler(db, 24*time.Hour)
	mymux.Handle("/webhook", &server.Handler{DB: db})
	err = http.ListenAndServe(os.Getenv("PORT"), server.ReqLog(mymux))
	if err != nil {
		log.Fatal("Server error", err)
	}
}