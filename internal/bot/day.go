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
		text := fmt.Sprintf("*%s* as *%s* | %s\n%s %s\n%s x %s\n%s\n\n",
			game.Date.Format("2/01"),
			game.Date.Format("15:04"),
			game.Place,
			game.Stage,
			game.Group,
			game.Team1,
			game.Team2,
			game.Location,
		)

		msg.Text += text
	}

	msg.ParseMode = "markdown"

	if _, err := c.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
