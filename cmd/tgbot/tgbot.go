package main

import "github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot"

func main() {
	app := tgbot.New(
		tgbot.WithConfig(),
	)

	app.Run()
}
