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

const formalizeServer = "http://formalize:8081"

func NewFormalizer() (*Formalizer, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 0
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

func (client *Formalizer) Formalize(trs string) (string, error) {
	res, err := client.TrsFormalizeWithResponse(context.TODO(), TrsFormalizeJSONRequestBody{
		Trs: trs,
	})
	if err != nil {
		return "", err
	}
	if res.StatusCode() != http.StatusOK {
		slog.Error("error requesting Formalize", "code", res.StatusCode())
		return "", errors.New("error requesting Formalize")
	}

	return res.JSON200.FormalTrs, nil
}

func (client *Formalizer) FixFormalized(trs string, formalTrs string, errorStr string) (string, error) {
	res, err := client.TrsFixWithResponse(context.TODO(), TrsFixJSONRequestBody{
		Error:     errorStr,
		FormalTrs: formalTrs,
		Trs:       trs,
	})
	if err != nil {
		return "", err
	}
	if res.StatusCode() != http.StatusOK {
		slog.Error("error requesting Formalize", "code", res.StatusCode())
		return "", errors.New("error requesting Formalize")
	}

	return res.JSON200.FormalTrs, nil
}
