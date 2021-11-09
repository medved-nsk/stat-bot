package main

import (
	"net/http"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wsgoggway/statbot/pkg/bot"
	"github.com/wsgoggway/statbot/pkg/initial"
	"github.com/wsgoggway/statbot/pkg/rest"
)

var (
	wg sync.WaitGroup
)

func main() {
	log, db, updates, telegramBot, userList := initial.InitProgram()

	wg.Add(2)

	go func() {
		h := rest.APIHandler{DB: db}
		e := echo.New()
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}))
		e.Static("/stat", "page")
		e.GET("/api", h.GetStat)

		e.Logger.Fatal(e.Start(":4000"))
	}()

	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.Chat.ID > 0 {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true

			bothandler := bot.BotHandler{DB: db, Message: msg, UserList: userList}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "stat":
					bothandler.GetStatCommand("all")
				case "daily":
					bothandler.GetStatCommand("daily")
				case "last7days":
					bothandler.GetStatCommand("seven")
				default:
					continue
				}

				if _, err := telegramBot.Send(bothandler.Message); err != nil {
					log.Panic(err)
				}
				continue
			}

			bothandler.StoreMessage(update)
		}
	}()

	wg.Wait()
}
