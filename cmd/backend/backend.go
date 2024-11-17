package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

func cli() {
	r := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(r)
	if err != nil {
		ExitWithError(err)
	}

	model, err := mclient.NewMistralClient()
	if err != nil {
		ExitWithError(err)
	}

	slog.Info("Executing model request")

	answers, err := usecases.AskKnowledgeBase(context.Background(), model, string(data))
	if err != nil {
		ExitWithError(err)
	}

	for _, answer := range answers {
		fmt.Printf("model=%s question=%s\ncontext=%v", answer.Model, answer.Answer, answer.Context)
	}
}

func main() {
	useCli := flag.Bool("cli", false, "run with cli interface")

	flag.Parse()

	if *useCli {
		cli()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := tgbot.New(
		ctx,
		tgbot.WithConfig(),
	)
	if err != nil {
		ExitWithError(err)
	}

	app.Run(ctx)
}

func ExitWithError(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
