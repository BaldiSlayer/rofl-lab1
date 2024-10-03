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
	slog.Info("start parsing with input:\n"+input)

	if input == "" {
		return nil, &ParseError{
			LlmMessage: "система должна содержать хотя бы одно правило переписывания и его интерпретацию",
			Summary:    "empty input",
		}
	}

	trs, err := p.parse(input)
	if err != nil {
		return nil, toParseError(err)
	}
	return trs, err
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
		Rules:           ToRulesDTO(trs.Rules, trs.Variables),
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

func ToRulesDTO(rules []rulesparser.Rule, variables []models.Lexem) []Rule {
	variablesMap := toMap(variables)

	res := []Rule{}
	for _, el := range rules {
		res = append(res, Rule{
			Lhs: toSubexpressionDTO(el.Lhs, variablesMap),
			Rhs: toSubexpressionDTO(el.Rhs, variablesMap),
		})
	}
	return res
}

func toSubexpressionDTO(subexpr rulesparser.Subexpression, variables map[string]struct{}) Subexpression {
	args := []interface{}{}
	if subexpr.Args != nil {
		for _, el := range *subexpr.Args {
			arg := toSubexpressionDTO(el, variables)
			args = append(args, arg)
		}
	}

	name := subexpr.Letter.String()
	_, isVariable := variables[name]

	return Subexpression{
		Args:   &args,
		Letter: Letter{
			IsVariable: isVariable,
			Name:       name,
		},
	}
}

func toInterpretsDTO(interprets []trsinterprets.Interpretation) []Interpretation {
	res := []Interpretation{}
	for _, el := range interprets {
		res = append(res, Interpretation{
			Args:      el.Args,
			Monomials: toMonomialsDTO(el.Monomials),
			Name:      el.Name,
		})
	}
	return res
}

func toMonomialsDTO(monomials []trsinterprets.Monomial) []Monomial {
	res := []Monomial{}
	for _, el := range monomials {
		if el.Constant != nil {
			res = append(res, NewConstantMonomial(*el.Constant))
			continue
		}
		res = append(res, NewProductMonomial(toFactorsDTO(*el.Factors)))
	}
	return res
}

func toFactorsDTO(factors []trsinterprets.Factor) []Factor {
	res := []Factor{}
	for _, tmp := range factors {
		el := tmp
		factor := Factor{
			Coefficient: nil,
			Power:       nil,
			Variable:    el.Variable,
		}
		if el.Coefficient != 1 {
			factor.Coefficient = &el.Coefficient
		}
		if el.Power != 1 {
			factor.Power = &el.Power
		}
		res = append(res, factor)
	}
	return res
}

func toVariablesDTO(variables []models.Lexem) []string {
	res := []string{}
	for _, el := range variables {
		res = append(res, el.String())
	}
	return res
}

func toMap(slice []models.Lexem) map[string]struct{} {
	res := make(map[string]struct{}, len(slice))
	for _, el := range slice {
		res[el.String()] = struct{}{}
	}
	return res
}
