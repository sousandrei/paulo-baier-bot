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
		msg.Text += game.String()
	}

	msg.ParseMode = "markdown"

	if _, err := c.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
