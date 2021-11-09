package rest

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/wsgoggway/statbot/pkg/initial"
	"github.com/wsgoggway/statbot/pkg/structs"
)

var (
	log = initial.Log
)

type APIHandler struct {
	DB *sqlx.DB
}

func (h *APIHandler) GetStat(c echo.Context) error {
	var data []structs.Stat
	row, _ := h.DB.Queryx(`SELECT
	count() AS count,
	from_username AS username,
	from_first_name AS first_name,
	from_last_name AS last_name,
	sum(char_length(text)) AS length_of_write_text
FROM telegramstat
GROUP BY username, first_name, last_name
ORDER BY
	count DESC,
	length_of_write_text DESC`)

	for row.Next() {
		var _data structs.Stat
		err := row.StructScan(&_data)
		if err != nil {
			log.Info(err)
		}
		data = append(data, _data)
	}
	return c.JSON(http.StatusOK, data)
}
