package formalizeclient

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../../docs/formalize-api.yaml

type Formalizer struct {
	*ClientWithResponses
}

// TODO вынести в конфиг, хардкодить неудобно
const (
	formalizeServer = "http://formalize:8000"
	retryMax        = 2
)

func NewFormalizer() (*Formalizer, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retryMax
	standardClient := retryClient.StandardClient() // *http.Client

	c, err := NewClientWithResponses(formalizeServer, WithHTTPClient(standardClient))
	if err != nil {
		return nil, err
	}

	slog.Info("initialized formalize client")

	return &Formalizer{
		ClientWithResponses: c,
	}, nil
}

type FormalizeResultDTO struct {
	FormalizedTrs    string
	ErrorDescription *string
}

func (client *Formalizer) Formalize(ctx context.Context, trs string) (FormalizeResultDTO, error) {
	res, err := client.TrsFormalizeWithResponse(ctx, TrsFormalizeJSONRequestBody{
		Trs: trs,
	})
	if err != nil {
		return FormalizeResultDTO{}, err
	}
	if res.StatusCode() != http.StatusOK {
		slog.Error("error requesting Formalize", "code", res.StatusCode())
		return FormalizeResultDTO{}, errors.New("error requesting Formalize")
	}

	return FormalizeResultDTO{
		FormalizedTrs:    res.JSON200.FormalTrs,
		ErrorDescription: res.JSON200.Error,
	}, nil
}

func (client *Formalizer) FixFormalized(ctx context.Context, trs string, formalTrs string, errorStr string) (FormalizeResultDTO, error) {
	res, err := client.TrsFixWithResponse(ctx, TrsFixJSONRequestBody{
		Error:     errorStr,
		FormalTrs: formalTrs,
		Trs:       trs,
	})
	if err != nil {
		return FormalizeResultDTO{}, err
	}
	if res.StatusCode() != http.StatusOK {
		slog.Error("error requesting Formalize", "code", res.StatusCode())
		return FormalizeResultDTO{}, errors.New("error requesting Formalize")
	}

	return FormalizeResultDTO{
		FormalizedTrs:    res.JSON200.FormalTrs,
		ErrorDescription: res.JSON200.Error,
	}, nil
}
