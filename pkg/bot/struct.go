package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/wsgoggway/statbot/pkg/structs"
)

type BotHandler struct {
	DB       *sqlx.DB
	Message  tgbotapi.MessageConfig
	UserList map[int64]structs.User
}
