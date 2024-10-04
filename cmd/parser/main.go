package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(r)
	if err != nil {
		slog.Error("error while reading input string", "err", err)
		os.Exit(1)
	}

	trs, err := trsparser.Parser{}.Parse(string(input))
	if err != nil {
		slog.Error("error while parsing input string", "err", err)
		os.Exit(1)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(trs)
	if err != nil {
		slog.Error("error encoding trs", "err", err)
		os.Exit(1)
	}
}
