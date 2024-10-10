package trsparser

import (
	"fmt"
	"strconv"
	"strings"

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
	res := []string{}
	for _, variable := range variables {
		res = append(res, variable.String())
	}
	return res
}

func toInterpretsDTO(interprets []trsinterprets.Interpretation) []string {
	res := []string{}
	for _, interpret := range interprets {
		lhs := getLhsString(interpret)
		rhs := toInterpretationDTO(interpret.Monomials)
		res = append(res, fmt.Sprintf("%s=%s", lhs, rhs))
	}
	return res
}

func toRulesDTO(rules []rulesparser.Rule) []string {
	res := []string{}
	for _, rule := range rules {
		res = append(res, toRuleDTO(rule))
	}
	return res
}

func toRuleDTO(rule rulesparser.Rule) string {
	return fmt.Sprintf("%s=%s", toSubexpressionDTO(rule.Lhs), toSubexpressionDTO(rule.Rhs))
}

func toSubexpressionDTO(subexpr rulesparser.Subexpression) string {
	if subexpr.Args == nil || len(*subexpr.Args) == 0 {
		return subexpr.Letter.String()
	}

	args := []string{}
	for _, el := range *subexpr.Args {
		args = append(args, toSubexpressionDTO(el))
	}

	return fmt.Sprintf("%s(%s)", subexpr.Letter.String(), strings.Join(args, ","))
}

func getLhsString(interpret trsinterprets.Interpretation) string {
	if len(interpret.Args) == 0 {
		return interpret.Name
	}
	args := strings.Join(interpret.Args, ",")
	return fmt.Sprintf("%s(%s)", interpret.Name, args)
}

func toInterpretationDTO(m []trsinterprets.Monomial) string {
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
