package bot

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocarina/gocsv"
)

// =================================================
// TODO: separate

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

// =================================================

type Client struct {
	bot   *tgbotapi.BotAPI
	games []Game
}

func New() *Client {
	//TODO: receive conf
	token := os.Getenv("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("Error creating bot", err)
	}

	fmt.Println("Authorized on account", bot.Self.UserName)

	allGames := []Game{}
	if err := gocsv.UnmarshalBytes(rawGames, &allGames); err != nil {
		fmt.Println("Error unmarshalling games.csv", err)
	}

	return &Client{
		bot:   bot,
		games: allGames,
	}
}

func (c *Client) Run() error {
	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = 60

	updates := c.bot.GetUpdatesChan(updConfig)

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
			c.handleNextCommand(&msg)
		}
	}

	return nil
}

func (c *Client) handleNextCommand(msg *tgbotapi.MessageConfig) {
	var game *Game
	for _, g := range c.games {
		if !time.Now().After(g.Date.Time) {
			game = &g
			break
		}
	}

	if game == nil {
		msg.Text = "No more games"
		if _, err := c.bot.Send(msg); err != nil {
			fmt.Println("Error sending message:", err)
		}

		return
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

	if _, err := c.bot.Send(msg); err != nil {
		fmt.Println("Error sending message:", err)
	}
}
