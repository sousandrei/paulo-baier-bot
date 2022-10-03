package bot

import (
	"fmt"

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
		}
	}

	return nil
}
