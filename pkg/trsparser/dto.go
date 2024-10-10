package trsparser

import (
	"strconv"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
	rulesparser "github.com/BaldiSlayer/rofl-lab1/internal/parser/trsrules"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../docs/trs-parser-api.yaml

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

func toVariablesDTO(variables []models.Lexem) []string {
	return nil
}

func toInterpretsDTO(interprets []trsinterprets.Interpretation) []string {
	return nil
}

func ToRulesDTO(rules []rulesparser.Rule, variables []models.Lexem) []string {
	return nil
}

func TranslateInterpretation(m []trsinterprets.Monomial) string {
	result := ""
	count := len(m)
	for i, monom := range m {
		if monom.Constant != nil {
			result += strconv.Itoa(*monom.Constant)
		} else if monom.Factors != nil {
			countMonoms := len(*monom.Factors)
			for j, factor := range *monom.Factors {
				if factor.Coefficient != 1 {
					result += strconv.Itoa(factor.Coefficient) + "*"
				}
				result += factor.Variable
				if factor.Power != 1 {
					result += "**" + strconv.Itoa(factor.Power)
				}
				if j != countMonoms-1 {
					result += "*"
				}
			}
		}
		if i != count-1 {
			result += "+"
		}
	}
	return result
}
