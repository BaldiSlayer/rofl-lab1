package usecases

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/formalizeclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/interpretclient"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

type TrsUseCases struct {
	parser trsparser.Parser
	interpret interpretclient.Interpreter
	formalize formalizeclient.Formalizer
}

func (uc TrsUseCases) ExtractFormalTrs(request string) (trsparser.Trs, string, error) {
	formalizedTrs, err := uc.formalize.Formalize(request)
	if err != nil {
		return trsparser.Trs{}, "", err
	}

	trs, err := uc.parser.Parse(formalizedTrs)
	if err != nil {
		return trsparser.Trs{}, formalizedTrs, err
	}

	return *trs, formalizedTrs, nil
}

func (uc TrsUseCases) FixFormalTrs(request, formalTrs string, parseError trsparser.ParseError) (trsparser.Trs, string, error) {
	// TODO: pass error
	formalizedTrs, err := uc.formalize.Formalize(request)
	if err != nil {
		return trsparser.Trs{}, "", err
	}

	trs, err := uc.parser.Parse(formalizedTrs)
	if err != nil {
		return trsparser.Trs{}, formalizedTrs, err
	}

	return *trs, formalizedTrs, nil
}

func (uc TrsUseCases) InterpretFormalTrs(trs trsparser.Trs) (string, error) {
	return uc.interpret.Interpret(trs)
}
