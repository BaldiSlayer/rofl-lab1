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
	formalizedTrs, err := uc.formalize.Formalize(ctx, request)
	for i := 0; i < retryMax && formalizedTrs.ErrorDescription != nil; i++ {
		slog.Info("got error from formalize", "error", *formalizedTrs.ErrorDescription)
		formalizedTrs, err = uc.formalize.FixFormalized(ctx, request, formalizedTrs.FormalizedTrs, *formalizedTrs.ErrorDescription)
	}
	if err != nil {
		return ExtractData{}, err
	}
	if formalizedTrs.ErrorDescription != nil {
		return ExtractData{}, fmt.Errorf("error formalizing trs: %s", *formalizedTrs.ErrorDescription)
	}

	trs, err := uc.parser.Parse(formalizedTrs)
	if err != nil {
		return ExtractData{
			Trs:           trsparser.Trs{},
			FormalizedTrs: formalizedTrs,
		}, err
	}

	return ExtractData{
		Trs:           *trs,
		FormalizedTrs: formalizedTrs,
	}, nil
}

// TODO: refactor duplicate code in ExtractFormalTrs and FixFormalTrs
func (uc *TrsUseCases) ExtractFormalTrs(ctx context.Context, request string) (ExtractData, error) {
	formalizedTrs, err := uc.formalize.Formalize(ctx, request)
	for i := 0; i < retryMax && result.ErrorDescription != nil; i++ {
		slog.Info("got error from formalize", "error", err)
		result, err = uc.formalize.FixFormalized(ctx, request, result.FormalizedTrs, *result.ErrorDescription)
	}
	if err != nil {
		return trsparser.Trs{}, "", err
	}
	if result.ErrorDescription != nil {
		return trsparser.Trs{}, "", fmt.Errorf("error formalizing trs: %s", *result.ErrorDescription)
	}

	trs, err := uc.parser.Parse(result.FormalizedTrs)
	if err != nil {
		return ExtractData{
			Trs:           trsparser.Trs{},
			FormalizedTrs: formalizedTrs,
		}, err
	}

	return ExtractData{
		Trs:           *trs,
		FormalizedTrs: formalizedTrs,
	}, nil
}

func (uc *TrsUseCases) FixFormalTrs(ctx context.Context, request, formalTrs, errorDescription string) (ExtractData, error) {
	formalizedTrs, err := uc.formalize.FixFormalized(ctx, request, formalTrs, errorDescription)
	if err != nil {
		return ExtractData{}, err
	}

	trs, err := uc.parser.Parse(formalizedTrs)
	if err != nil {
		return ExtractData{
			Trs:           trsparser.Trs{},
			FormalizedTrs: formalizedTrs,
		}, err
	}

	return ExtractData{
		Trs:           *trs,
		FormalizedTrs: formalizedTrs,
	}, nil
}


func (uc *TrsUseCases) FixFormalTrs(ctx context.Context, request, formalTrs, errorDescription string) (ExtractData, error) {
	formalizedTrs, err := uc.formalize.FixFormalized(ctx, request, formalTrs, errorDescription)
	if formalizedTrs.ErrorDescription != nil {
		return ExtractData{}, fmt.Errorf("error formalizing trs: %s", *formalizedTrs.ErrorDescription)
	}
	for i := 0; i < retryMax && formalizedTrs.ErrorDescription != nil; i++ {
		slog.Info("got error from formalize", "error", err)
		formalizedTrs, err = uc.formalize.FixFormalized(ctx, request, formalizedTrs.FormalizedTrs, *formalizedTrs.ErrorDescription)
	}
	if err != nil {
		return ExtractData{}, err
	}
	if formalizedTrs.ErrorDescription != nil {
		return ExtractData{}, fmt.Errorf("error formalizing trs: %s", *formalizedTrs.ErrorDescription)
	}

	trs, err := uc.parser.Parse(formalizedTrs)
	if err != nil {
		return ExtractData{
			Trs:           trsparser.Trs{},
			FormalizedTrs: formalizedTrs,
		}, err
	}

	return ExtractData{
		Trs:           *trs,
		FormalizedTrs: formalizedTrs,
	}, nil
}

func (uc *TrsUseCases) InterpretFormalTrs(ctx context.Context, trs trsparser.Trs) (string, error) {
	return uc.interpret.Interpret(ctx, trs)
}
