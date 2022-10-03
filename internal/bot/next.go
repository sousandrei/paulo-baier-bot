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
