package initial

import (
	"os"
	"strings"

	"github.com/wsgoggway/statbot/pkg/structs"

	_ "github.com/ClickHouse/clickhouse-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

var (
	Log *zap.SugaredLogger
)

// InitProgram ...
func InitProgram() (*zap.SugaredLogger, *sqlx.DB, tgbotapi.UpdatesChannel, *tgbotapi.BotAPI, map[int64]structs.User) {
	logger, _ := zap.NewDevelopment()
	Log := logger.Sugar()

	err := godotenv.Load()
	if err != nil {
		Log.Warn("Error loading .env file")
	}

	connect, err := sqlx.Open("clickhouse", os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		Log.Fatal(err)
	}

	_, err = connect.Exec(`
	CREATE TABLE IF NOT EXISTS telegramstat (
		update_id Int64,
		message_id Int64,
		from_id Int64,
		from_first_name String,
		from_last_name String,
		from_username String,
		chat_id Int64, 
		chat_title String, 
		chat_type String,
		date DateTime,
		text String
	)
	ENGINE = MergeTree()
	PARTITION BY (toYYYYMMDD(date))
	ORDER BY (from_id, date)
`)
	if err != nil {
		Log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
	} else {
		bot.Debug = false
	}
	log.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	users := make(map[int64]structs.User)

	row, _ := connect.Queryx("SELECT from_id FROM telegramstat GROUP BY from_id")
	for row.Next() {
		var _data struct {
			ID int64 `db:"from_id"`
		}
		err := row.StructScan(&_data)
		if err != nil {
			log.Info(err)
		}

		userRow, err := connect.Queryx(`SELECT
    from_first_name,
    from_last_name,
    from_username
FROM telegramstat
WHERE from_id = ? AND ((from_username != '') OR (from_first_name != '') OR (from_last_name != ''))
ORDER BY date DESC
LIMIT 1`, _data.ID)
		if err != nil {
			log.Info(err)
		}
		for userRow.Next() {
			var user = structs.User{}
			err := userRow.StructScan(&user)
			if err != nil {
				log.Info(err)
			}
			if strings.HasPrefix(user.Username, "unkwn<") || strings.HasSuffix(user.Username, "bot") {
				continue
			}

			users[_data.ID] = user
		}
	}

	return Log, connect, updates, bot, users
}
