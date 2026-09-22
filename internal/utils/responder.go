package utils

import ("database/sql"
"net/http"
"encoding/json"
"bytes"
)
// The Responder

var Categories InlineKeyboardMarkup = InlineKeyboardMarkup{
	InlineKeyboard: [][]InlineKeyboardButton{
		{
			{Text: "Go", CallbackData: "go"},
			{Text: "Python", CallbackData: "python"},
			{Text: "Java", CallbackData: "java"},
			{Text: "DevOps", CallbackData: "devops"},
		},
	},
}

var Locations InlineKeyboardMarkup = InlineKeyboardMarkup{
	InlineKeyboard: [][]InlineKeyboardButton{
		{
			{Text: "Warszawa", CallbackData: "Warszawa"},
			{Text: "Kraków", CallbackData: "Kraków"},
			{Text: "Wrocław", CallbackData: "Wrocław"},
			{Text: "Poznań", CallbackData: "Poznań"},
			{Text: "Gdańsk", CallbackData: "Gdańsk"},
			{Text: "Praha", CallbackData: "Praha"},
		},
	},
}

var Workplaces InlineKeyboardMarkup = InlineKeyboardMarkup{
	InlineKeyboard: [][]InlineKeyboardButton{
		{
			{Text: "Remote", CallbackData: "remote"},
			{Text: "Office", CallbackData: "office"},
			{Text: "Hybrid", CallbackData: "hybrid"},
		},
	},
}

var ExperienceLevels InlineKeyboardMarkup = InlineKeyboardMarkup{
	InlineKeyboard: [][]InlineKeyboardButton{
		{
			{Text: "Junior", CallbackData: "junior"},
			{Text: "Mid", CallbackData: "mid"},
			{Text: "Senior", CallbackData: "senior"},
			{Text: "Manager", CallbackData: "manager"},
		},
	},
}

type Responder struct {
	ChatID int64
	UserID int64
}

type InlineKeyboardButton struct {
    Text         string `json:"text"`
    CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
    InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type SendMessageRequest struct {
    ChatID      int64                  `json:"chat_id"`
    Text        string                 `json:"text"`
    ReplyMarkup *InlineKeyboardMarkup  `json:"reply_markup,omitempty"`
}

func (r *Responder) Sendmessage(m string) error {
	var message SendMessageRequest
	message.ChatID = r.ChatID
	message.Text = m
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}

func (r *Responder) StartSubcription(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO pending_subscriptions (userid, chatid, step) VALUES ($1, $2, $3)", r.UserID, r.ChatID, "category")
	if err != nil {
		return err
	} 
	var message SendMessageRequest
	message.ChatID = r.ChatID
	message.Text = "Please select a category:"
	message.ReplyMarkup = &Categories
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}


func (r *Responder) SendCallback(cb string, db *sql.DB) error {
	resp, err := db.Query("SELECT step FROM pending_subscriptions WHERE userid = $1 AND chatid = $2", r.UserID, r.ChatID)
	if err != nil {
		return err
	}
	defer resp.Close()
	var step string
	for resp.Next() {
		err = resp.Scan(&step)
		if err != nil {
			return err
		}
	}
	switch step {
	case "category":
		_, err = db.Exec("UPDATE pending_subscriptions SET category = $1, step = $2 WHERE userid = $3 AND chatid = $4", cb, "location", r.UserID, r.ChatID)
		if err != nil {
			return err
		}
		var message SendMessageRequest
		message.ChatID = r.ChatID
		message.Text = "Please select a location:"
		message.ReplyMarkup = &Locations
		jsonData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
	case "location":
		_, err = db.Exec("UPDATE pending_subscriptions SET location = $1, step = $2 WHERE userid = $3 AND chatid = $4", cb, "workplace", r.UserID, r.ChatID)
		if err != nil {
			return err
		}
		var message SendMessageRequest
		message.ChatID = r.ChatID
		message.Text = "Please select a workplace type:"
		message.ReplyMarkup = &Workplaces
		jsonData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
	case "workplace":
		_, err = db.Exec("UPDATE pending_subscriptions SET workplace = $1, step = $2 WHERE userid = $3 AND chatid = $4", cb, "experience", r.UserID, r.ChatID)
		if err != nil {
			return err
		}
		var message SendMessageRequest
		message.ChatID = r.ChatID
		message.Text = "Please select an experience level:"
		message.ReplyMarkup = &ExperienceLevels
		jsonData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
	case "experience":
		_, err = db.Exec("UPDATE pending_subscriptions SET experience = $1, step = $2 WHERE userid = $3 AND chatid = $4", cb, "completed", r.UserID, r.ChatID)
		if err != nil {
			return err
		}
		var message SendMessageRequest
		message.ChatID = r.ChatID
		message.Text = "Subscription completed! You will now receive job offers based on your preferences."
		jsonData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		err = Finalize(r.ChatID, r.UserID, db)
		if err != nil {
			return err
		}
		_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
	default:
		var message SendMessageRequest
		message.ChatID = r.ChatID
		message.Text = "An error occurred. Please try again."
		jsonData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		_, err = http.Post("https://api.telegram.org/bot8949513390:AAGC-U8q7KzoKC8_VGpLZF4Nzczurk93m3I/sendMessage", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
	}
	return nil
}

func Finalize (chatid int64, userid int64, db *sql.DB) error {
	rows, err := db.Query("SELECT category, location, workplace, experience FROM pending_subscriptions WHERE chatid = $1 AND userid = $2 AND step = 'completed'", chatid, userid)
	if err != nil {
		return err
	}
	defer rows.Close()
	var row Subscription
	for rows.Next() {
		err = rows.Scan(&row.Category, &row.Location, &row.Workplace, &row.Experience)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("INSERT INTO subscriptions (userid, chatid, category, location, workplace, experience) VALUES ($1, $2, $3, $4, $5, $6)", userid, chatid, row.Category, row.Location, row.Workplace, row.Experience)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM pending_subscriptions WHERE chatid = $1 AND userid = $2", chatid, userid)
	if err != nil {
		return err
	}
	return nil
}