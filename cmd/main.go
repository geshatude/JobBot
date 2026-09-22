package main

import ("database/sql"
"log"
"github.com/geshatude/JobBot/internal/server"
_"github.com/lib/pq"
"github.com/geshatude/JobBot/internal/utils"
"net/http"
"time"
"golang.ngrok.com/ngrok/v2"
"context")

func connectNgrok() {
	fwd, err := ngrok.Forward(context.Background(),
		ngrok.WithUpstream("http://localhost:8085"),
		ngrok.WithURL("https://email-atonable-requisite.ngrok-free.dev"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Available at:", fwd.URL())
	select {}
}

func main() {
	mymux := http.NewServeMux()
	db, err := sql.Open("postgres", "postgres://postgres:mysecretpassword@127.0.0.1:5432/jobbot?sslmode=disable")
	if err != nil {
		log.Fatal("database connection error", err)
	}
	go utils.StartScheduler(db, 24*time.Hour)
	mymux.Handle("/telegram-webhook", &server.Handler{Database: db})
	go http.ListenAndServe(":8085", server.ReqLog(mymux))
	connectNgrok()
}