package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousandrei/paulobaierbot/internal/games"
)

func (c *Client) handleSearchTeamCommand(msg *tgbotapi.MessageConfig, str string) error {
	var games []games.Game
	for _, g := range c.games {
		if strings.Contains(strings.ToLower(g.Team1), strings.ToLower(str)) ||
			strings.Contains(strings.ToLower(g.Team2), strings.ToLower(str)) {
			games = append(games, g)
		}
	}

	if len(games) < 1 {
		msg.Text = "No games"
		if _, err := c.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

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
