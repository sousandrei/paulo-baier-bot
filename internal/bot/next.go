package bot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousandrei/paulobaierbot/internal/games"
)

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

	msg.Text = game.String()
	msg.ParseMode = "markdown"

	if _, err := c.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
