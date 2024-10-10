package trsparser

import (
	"log/slog"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/lexer"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
	rulesparser "github.com/BaldiSlayer/rofl-lab1/internal/parser/trsrules"
)

type Parser struct{}

func (p Parser) Parse(input string) (*Trs, error) {
	slog.Info("start parsing")

	if input == "" {
		return nil, &ParseError{
			LlmMessage: "система должна содержать хотя бы одно правило переписывания и его интерпретацию",
			Summary:    "empty input",
		}
	}

	trs, err := p.parse(input)
	return trs, toParseErrorDTO(err)
}

func (p Parser) parse(input string) (*Trs, error) {
	slog.Info("run lexer")

	l := lexer.Lexer{
		Text: input,
	}

	err := l.Process()
	if err != nil {
		return nil, err
	}

	slog.Info("run rules parser")

	trs, rest, err := rulesparser.ParseRules(l.Lexem)
	if err != nil {
		return nil, err
	}

	slog.Info("consume separators")

	rest, err = p.consumeSeparators(rest)
	if err != nil {
		return nil, err
	}

	slog.Info("run interprets parser")

	interpretsParser := trsinterprets.NewParser(trsinterprets.ToInputChannel(rest), trs.Constructors)

	interprets, err := interpretsParser.Parse()
	if err != nil {
		return nil, err
	}

	slog.Info("convert to DTO")

	return &Trs{
		Interpretations: toInterpretsDTO(interprets),
		Rules:           toRulesDTO(trs.Rules),
		Variables:       toVariablesDTO(trs.Variables),
	}, nil
}

func (p *Parser) consumeSeparators(lexems []models.Lexem) ([]models.Lexem, error) {
	lenBefore := len(lexems)
	for len(lexems) > 0 {
		lexem := lexems[0].Type()
		if lexem != models.LexSEPARATOR && lexem != models.LexEOL {
			break
		}
		lexems = lexems[1:]
	}
	if lenBefore == len(lexems) {
		return nil, &ParseError{
			LlmMessage: "не найдены разделители (-) между правилами и интерпретациями",
			Summary:    "no separators found",
		}
	}
	return lexems, nil
}
