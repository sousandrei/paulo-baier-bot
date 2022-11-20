package bot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousandrei/paulobaierbot/internal/games"
)

type Client struct {
	bot   *tgbotapi.BotAPI
	games []games.Game
}

func New(token string) (*Client, error) {
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

		case "today":
			today := time.Now().Format("2/01")
			err := c.handleDayCommand(&msg, today)
			if err != nil {
				return fmt.Errorf("error handling day command: %w", err)
			}

		case "day":
			txt := update.Message.CommandArguments()
			day, err := time.Parse("2/01", txt)
			if err != nil {
				msg.Text = "Invalid date format. Use dd/mm"
				if _, err := c.bot.Send(msg); err != nil {
					return fmt.Errorf("error sending message: %w", err)
				}

				continue
			}
			err = c.handleDayCommand(&msg, day.Format("2/01"))
			if err != nil {
				return fmt.Errorf("error handling day command: %w", err)
			}

		case "team":
			txt := update.Message.CommandArguments()
			err := c.handleSearchTeamCommand(&msg, txt)
			if err != nil {
				return fmt.Errorf("error handling search command: %w", err)
			}

		}

	}

	return nil
}
