package utils

import ("database/sql"
"github.com/geshatude/JobBot/internal/server"
"net/http"
"encoding/json"
"bytes"
"log"
)

//The Sorter

type Subscription struct {
	Category string
	Location string
	Workplace string
	Experience string
}
type Match struct {
    UserID     int64
    ChatID     int64
    Guid       string
    ApplyUrl   string
    Category   string
    Location   string
    Workplace  string
    Experience string
}

func GetMatches(db *sql.DB) ([]Match, error) {
    rows, err := db.Query("
        SELECT s.userid, s.chatid, o.guid, o.applyurl, o.category, o.location, o.workplace, o.experience
        FROM subscriptions s
        JOIN job_offers o
          ON s.category = o.category
          AND s.location = o.location
          AND s.workplace = o.workplace
          AND s.experience = o.experience
        WHERE NOT EXISTS (
          SELECT 1 FROM sent_offers sn
          WHERE sn.userid = s.userid AND sn.guid = o.guid
        )"
	)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var matches []Match
    for rows.Next() {
        var m Match
        err = rows.Scan(&m.UserID, &m.ChatID, &m.Guid, &m.ApplyUrl, &m.Category, &m.Location, &m.Workplace, &m.Experience)
        if err != nil {
            return nil, err
        }
        matches = append(matches, m)
    }
    return matches, nil
}
func Notify(matches []Match) error {
	for _, match := range matches {
    var message server.SendMessageRequest
    message.ChatID = match.ChatID
    message.Text = "New job offer found!\n\nCategory: " + match.Category + "\nLocation: " + match.Location + "\nWorkplace: " + match.Workplace + "\nExperience: " + match.Experience + "\n\nApply here: " + match.ApplyUrl

    jsonData, err := json.Marshal(message)
    if err != nil {
        log.Println("Error marshalling message:", err)
        return
    }

    _, err = http.Post("https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Println("Error sending message:", err)
        return
    }
    _, err = db.Exec("INSERT INTO sent_offers (userid, guid) VALUES ($1, $2)", match.UserID, match.Guid)
    if err != nil {
        log.Println("Error inserting sent offer:", err)
        return
    }   
    }
    return nil
}