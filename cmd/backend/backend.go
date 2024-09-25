package main

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend"
)

func main() {
	app := backend.New(
		backend.WithConfig(),
		backend.WithLogger(),
	)

	app.Run()
}
