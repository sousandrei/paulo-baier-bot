package main

import (
	"fmt"

	"github.com/sousandrei/paulobaierbot/internal/bot"
)

func main() {
	bot, err := bot.New()
	if err != nil {
		fmt.Println("Error creating bot", err)
	}

	if err := bot.Run(); err != nil {
		fmt.Println("Error running bot:", err)
	}
}
