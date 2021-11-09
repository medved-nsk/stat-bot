package bot

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/gommon/log"
)

func (h *BotHandler) StoreMessage(update tgbotapi.Update) {
	var (
		updateID      = update.UpdateID
		messageID     = update.Message.MessageID
		fromID        = update.Message.From.ID
		fromFirstName = update.Message.From.FirstName
		fromLastName  = update.Message.From.LastName
		fromUsername  = update.Message.From.UserName
		chatID        = update.Message.Chat.ID
		chatTitle     = update.Message.Chat.Title
		chatType      = update.Message.Chat.Type
		text          = update.Message.Text

		tx, _   = h.DB.Begin()
		stmt, _ = tx.Prepare("INSERT INTO telegramstat (update_id,message_id,from_id,from_first_name,from_last_name,from_username,chat_id,chat_title,chat_type,date,text) VALUES (?,?,?,?,?,?,?,?,?,toDateTime(?),?)")
	)
	defer stmt.Close()

	tm := time.Unix(int64(update.Message.Date), 0).Format("2006-01-02 15:04:05")
	if _, err := stmt.Exec(
		updateID,
		messageID,
		fromID,
		fromFirstName,
		fromLastName,
		fromUsername,
		chatID,
		chatTitle,
		chatType,
		tm,
		text,
	); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
