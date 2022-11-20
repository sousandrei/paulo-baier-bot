package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/sousandrei/paulobaierbot/internal/bot"
)

type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN"`
}

func main() {
	var cfg Config
	err := envconfig.Process("myapp", &cfg)
	if err != nil {
		fmt.Println("error processing envconfig: %w", err)
		return
	}

	bot, err := bot.New(cfg.TelegramToken)
	if err != nil {
		fmt.Println("Error creating bot", err)
		return
	}

	if err := bot.Run(); err != nil {
		fmt.Println("Error running bot:", err)
		return
	}
}
