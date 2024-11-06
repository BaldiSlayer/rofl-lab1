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
		slog.Error(err.Error())
		os.Exit(1)
	}

	model, err := mclient.NewMistralClient()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	qa, err := usecases.LoadQABase()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = model.InitContext(context.Background(), qa)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Executing model request")

	answer, err := usecases.AskKnowledgeBase(context.Background(), model, string(data))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(answer)
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
		slog.Error(err.Error())
		os.Exit(1)
	}

	app.Run(ctx)
}
