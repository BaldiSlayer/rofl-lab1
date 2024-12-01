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
		ExitWithError("error reading request from stdin", "error", err)
	}

	model, err := mclient.NewMistralClient()
	if err != nil {
		ExitWithError("failed to init llm client", "error", err)
	}

	slog.Info("Executing model request")

	askResults, err := usecases.AskKnowledgeBase(context.Background(), model, string(data))
	if err != nil {
		ExitWithError("error requesting knowledge base", "error", err)
	}

	for _, answer := range askResults.Answers {
		if !answer.UseContext {
			fmt.Printf(
				"model=%s question=%s\n",
				answer.Model,
				answer.Answer,
			)

			continue
		}

		fmt.Printf(
			"model=%s question=%s\ncontext=%v",
			answer.Model,
			answer.Answer,
			askResults.QuestionsContext,
		)
	}
}

func main() {
	useCli := flag.Bool("cli", false, "run with cli interface")
	callbackMode := flag.Bool("callback-mode", false, "run tg bot in long polling mode")

	flag.Parse()

	if *useCli {
		cli()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := tgbot.New(
		ctx,
		*callbackMode,
		tgbot.WithConfig(),
	)
	if err != nil {
		ExitWithError("failed to init telegram client", "error", err.Error())
	}

	app.Run(ctx)
}

func ExitWithError(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
