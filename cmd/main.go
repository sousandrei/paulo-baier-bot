package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocarina/gocsv"
)

//go:embed games.csv
var rawGames []byte

type Game struct {
	Date     GameTime `csv:"Date"`
	Team1    string   `csv:"Team 1"`
	Match    string   `csv:"Match"`
	Team2    string   `csv:"Team 2"`
	Group    string   `csv:"Group"`
	Location string   `csv:"Location"`
	Stage    string   `csv:"Stage"`
	Place    string   `csv:"Place"`
}

type GameTime struct {
	time.Time
}

func (date *GameTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2-Jan-06 3:04 PM", csv)
	return err
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")

	allGames := []Game{}
	if err := gocsv.UnmarshalBytes(rawGames, &allGames); err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		switch update.Message.Command() {
		case "next":
			var game *Game
			for _, g := range allGames {
				if !time.Now().After(g.Date.Time) {
					game = &g
					break
				}
			}

			if game == nil {
				msg.Text = "No more games"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

				continue
			}

			text := fmt.Sprintf("*%s* as *%s* | %s\n%s %s\n%s x %s\n%s",
				game.Date.Format("2/01"),
				game.Date.Format("15:04"),
				game.Place,
				game.Stage,
				game.Group,
				game.Team1,
				game.Team2,
				game.Location,
			)

			msg.Text = text
			msg.ParseMode = "markdown"

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

	}

}
