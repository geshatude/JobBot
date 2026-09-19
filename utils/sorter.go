package utils

import ("database/sql"
"github.com/geshatude/JobBot/internal")

type Subscription struct {
	Category string
	Location string
	Workplace string
	Experience string
}

func Match(db *sql.DB) error {
	rows, err := db.
}

func Notify() {}