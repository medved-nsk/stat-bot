package bot

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/wsgoggway/statbot/pkg/structs"
)

var (
	selectStat = `SELECT
	from_id,
	count() AS count,
	sum(char_length(text)) AS length_of_write_text
FROM telegramstat
GROUP BY from_id
ORDER BY
	count DESC,
	length_of_write_text DESC`

	selectDailyStat = `SELECT
	from_id,
	count() AS count,
	sum(char_length(text)) AS length_of_write_text
FROM telegramstat
WHERE date BETWEEN ? AND ?
GROUP BY from_id
ORDER BY
	count DESC,
	length_of_write_text DESC`
)

// GetStatCommand hanler for /stat command
func (bot *BotHandler) GetStatCommand(typeOfStat string) {
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	nextDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(time.Duration(24) * time.Hour).Unix()
	currentDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
	week := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(time.Duration(-24) * (7 * time.Hour)).Unix()

	var (
		row           *sqlx.Rows
		err           error
		totalMessages int64
		totalChar     int64
	)
	switch typeOfStat {
	case "daily":
		row, err = bot.DB.Queryx(selectDailyStat, currentDay, nextDay)
	case "seven":
		row, err = bot.DB.Queryx(selectDailyStat, week, currentDay)
	default:
		row, err = bot.DB.Queryx(selectStat)
	}
	if err != nil {
		log.Warn(err)
	}

	var line bytes.Buffer
	w := tabwriter.NewWriter(&line, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Сообщ.\tСимв.\tПользователь")

	for row.Next() {
		var _data structs.Stat
		err := row.StructScan(&_data)
		if err != nil {
			log.Info(err)
		}

		if user, ok := bot.UserList[_data.FromID]; ok {
			_data.Username = user.Username
			_data.FirstName = user.FirstName
			_data.LastName = user.LastName
		} else {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t", humanize.Comma(_data.Count), humanize.Comma(_data.LengthOfWriteText))
		if _data.FirstName != "" {
			fmt.Fprintf(w, "%s %s\n", _data.FirstName, _data.LastName)
		} else {
			fmt.Fprintf(w, "%s\n", _data.Username)
		}

		totalMessages += _data.Count
		totalChar += _data.LengthOfWriteText
	}
	w.Flush()

	bot.Message.Text = fmt.Sprintf(
		"```%s```\nВсего сообщений: %s (символов: %s)",
		line.String(),
		humanize.Comma(totalMessages),
		humanize.Comma(totalChar),
	)
}
