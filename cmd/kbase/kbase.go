package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(r)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	model, err := mclient.NewMistralClient()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Executing model request")

	answer, err := model.Ask(string(data))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(answer)
}
