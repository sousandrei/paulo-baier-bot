package bot

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousandrei/paulobaierbot/internal/games"
)

type Client struct {
	bot   *tgbotapi.BotAPI
	games []games.Game
}

func New() (*Client, error) {
	//TODO: receive conf
	token := os.Getenv("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot from api: %w", err)
	}

	// TODO: better logging
	fmt.Println("Authorized on account", bot.Self.UserName)

	games, err := games.GetGames()
	if err != nil {
		return nil, fmt.Errorf("error getting games: %w", err)
	}

	return &Client{
		bot:   bot,
		games: games,
	}, nil
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
			err := c.handleNextCommand(&msg)
			if err != nil {
				return fmt.Errorf("error handling next command: %w", err)
			}
		}
	}

	return nil
}

func (c *Client) handleNextCommand(msg *tgbotapi.MessageConfig) error {
	var game *games.Game
	for _, g := range c.games {
		if !time.Now().After(g.Date.Time) {
			game = &g
			break
		}
	}

	if game == nil {
		msg.Text = "No more games"
		if _, err := c.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		// TODO: something when it's over??
		return nil
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
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
