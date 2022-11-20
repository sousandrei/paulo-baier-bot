package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousandrei/paulobaierbot/internal/games"
)

func (c *Client) handleDayCommand(msg *tgbotapi.MessageConfig, day string) error {
	var games []games.Game
	for _, g := range c.games {
		if g.Date.Format("2/01") == day {
			games = append(games, g)
		}
	}

	if len(games) < 1 {
		msg.Text = "No games"
		if _, err := c.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		// TODO: something when it's over??
		return nil
	}

	for _, game := range games {
		msg.Text += game.String()
	}

	msg.ParseMode = "markdown"

	if _, err := c.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
