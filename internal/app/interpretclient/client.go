package interpretclient

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../../docs/interpret-api.yaml

type Interpreter struct {
	*ClientWithResponses
}

const interpretServer = "http://interpret:8081"

func NewInterpreter() (*Interpreter, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient() // *http.Client

	c, err := NewClientWithResponses(interpretServer, WithHTTPClient(standardClient))
	if err != nil {
		return nil, err
	}

	slog.Info("initialized interpret client")

	return &Interpreter{
		ClientWithResponses: c,
	}, nil
}

func (i *Interpreter) Interpret(ctx context.Context, trs trsparser.Trs) (string, error) {
	res, err := i.TrsInterpretWithResponse(ctx, TrsInterpretJSONRequestBody{
		Interpretations: trs.Interpretations,
		Rules:           trs.Rules,
		Variables:       trs.Variables,
	})
	if err != nil {
		return "", err
	}
	if res.StatusCode() != http.StatusOK {
		slog.Error("error requesting Interpret", "code", res.StatusCode())
		return "", errors.New("error requesting Interpret")
	}

	return res.JSON200.Answer, nil
}
