package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/formalizeclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/interpretclient"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

const (
	retryMax = 5
)

type TrsUseCases struct {
	parser    *trsparser.Parser
	interpret *interpretclient.Interpreter
	formalize *formalizeclient.Formalizer
}

func New() (*TrsUseCases, error) {
	interpreter, err := interpretclient.NewInterpreter()
	if err != nil {
		return nil, err
	}

	formalizer, err := formalizeclient.NewFormalizer()
	if err != nil {
		return nil, err
	}

	return &TrsUseCases{
		parser:    trsparser.NewParser(),
		interpret: interpreter,
		formalize: formalizer,
	}, nil
}

type ExtractData struct {
	Trs           trsparser.Trs
	FormalizedTrs string
}

func (uc *TrsUseCases) ExtractFormalTrs(ctx context.Context, request string) (ExtractData, error) {
	result, err := uc.formalize.Formalize(ctx, request)
	return uc.fixFormalTrs(ctx, request, result, err)
}

func (uc *TrsUseCases) FixFormalTrs(ctx context.Context, request, formalTrs, errorDescription string) (ExtractData, error) {
	result, err := uc.formalize.FixFormalized(ctx, request, formalTrs, errorDescription)
	return uc.fixFormalTrs(ctx, request, result, err)
}

func (uc *TrsUseCases) fixFormalTrs(ctx context.Context, request string, result formalizeclient.FormalizeResultDTO, err error) (ExtractData, error) {
	if result.ErrorDescription != nil {
		return ExtractData{}, fmt.Errorf("error formalizing trs: %s", *result.ErrorDescription)
	}
	for i := 0; i < retryMax && result.ErrorDescription != nil; i++ {
		slog.Info("got error from formalize", "error", err)
		result, err = uc.formalize.FixFormalized(ctx, request, result.FormalizedTrs, *result.ErrorDescription)
	}
	if err != nil {
		return ExtractData{}, err
	}
	if result.ErrorDescription != nil {
		return ExtractData{}, fmt.Errorf("error formalizing trs: %s", *result.ErrorDescription)
	}

	trs, err := uc.parser.Parse(result.FormalizedTrs)
	if err != nil {
		return ExtractData{
			Trs:           trsparser.Trs{},
			FormalizedTrs: result.FormalizedTrs,
		}, err
	}

	return ExtractData{
		Trs:           *trs,
		FormalizedTrs: result.FormalizedTrs,
	}, nil
}

func (uc *TrsUseCases) InterpretFormalTrs(ctx context.Context, trs trsparser.Trs) (string, error) {
	return uc.interpret.Interpret(ctx, trs)
}
