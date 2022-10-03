package main

import (
	"fmt"

	"github.com/sousandrei/paulobaierbot/internal/bot"
)

func main() {
	bot := bot.New()

	if err := bot.Run(); err != nil {
		fmt.Println("Error running bot:", err)
	}
}
