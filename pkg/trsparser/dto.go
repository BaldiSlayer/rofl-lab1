package trsparser

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
	rulesparser "github.com/BaldiSlayer/rofl-lab1/internal/parser/trsrules"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../docs/trs-parser-api.yaml

func NewConstantMonomial(constant int) Monomial {
	monomial := Monomial{}
	monomial.FromConstantMonomial(ConstantMonomial{
		Constant: constant,
	})
	return monomial
}

func NewProductMonomial(factors []Factor) Monomial {
	monomial := Monomial{}
	monomial.FromProductMonomial(ProductMonomial{
		Factors: factors,
	})
	return monomial
}

func NewSubexpression(sexpr Subexpression) interface{} {
	return sexpr
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
		Args: &args,
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

func toParseErrorDTO(err error) error {
	if err == nil {
		return nil
	}

	if parseError, ok := err.(*models.ParseError); ok {
		return &ParseError{
			LlmMessage: parseError.LlmMessage,
			Summary:    parseError.Message,
		}
	}

	return &ParseError{
		LlmMessage: "неизвестная ошибка",
		Summary:    err.Error(),
	}
}
