package server

import ("net/http"
"encoding/json"
"io"
"database/sql"
"log/slog"
"time"
"github.com/geshatude/JobBot/internal/utils")

// The Handler

type Handler struct {
	Database *sql.DB
}

type Update struct {
    UpdateID int      `json:"update_id"`
    Message  *Message `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type Message struct {
    MessageID int    `json:"message_id"`
    From      User   `json:"from"`
    Chat      Chat   `json:"chat"`
    Date      int    `json:"date"`
    Text      string `json:"text"`
}

type User struct {
    ID        int64  `json:"id"`
    IsBot     bool   `json:"is_bot"`
    FirstName string `json:"first_name"`
    Username  string `json:"username"`
}

type Chat struct {
    ID   int64  `json:"id"`
    Type string `json:"type"`
}

type CallbackQuery struct {
    ID      string  `json:"id"`
    From    User    `json:"from"`
    Message *Message `json:"message"`  
    Data    string  `json:"data"`      
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) error {
    update, err := io.ReadAll(r.Body)
    if err != nil {
        return err
    }
    var currentUpdate Update
    err = json.Unmarshal(update, &currentUpdate)
    if err != nil {
        return err
    }

    if currentUpdate.Message != nil {
        responder := &utils.Responder{
            ChatID: currentUpdate.Message.Chat.ID,
            UserID: currentUpdate.Message.From.ID,
        }
        if currentUpdate.Message.Text == "/start" {
            err = responder.Sendmessage("Hello! Welcome to JobBot. Please use the /subscribe command to begin.")
			if err != nil {
				return err
			}
        } else if currentUpdate.Message.Text == "/subscribe" {
            err = responder.StartSubcription(h.Database)
			 if err != nil {
            return err
    		}
		}
    } else if currentUpdate.CallbackQuery != nil {
        responder := &utils.Responder{
            ChatID: currentUpdate.CallbackQuery.Message.Chat.ID,
            UserID: currentUpdate.CallbackQuery.From.ID,
        }
        err = responder.SendCallback(currentUpdate.CallbackQuery.Data, h.Database)
        if err != nil {
            return err
        }
    }
    return nil
}

func ReqLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
		)
		reqAttrs := slog.Group("request", "ip", ip, "method", method, "url", url)
		slog.Info("request recieved", reqAttrs)

		next.ServeHTTP(w, r)

		slog.Info("request completed", "Duration", time.Since(start))
	})
}