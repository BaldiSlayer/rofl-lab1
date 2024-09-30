package interprets

import "fmt"

func checkMonomials(monomials []Monomial, args []string) *ParseError {
	c := monomialChecker{
		definedVars: toMap(args),
	}

	for _, el := range monomials {
		err := c.checkMonomial(el)
		if err != nil {
			return err
		}
	}

	return nil
}

type monomialChecker struct {
	definedVars map[string]struct{}
}

func (c *monomialChecker) checkMonomial(monomial Monomial) *ParseError {
	if _, ok := c.definedVars[monomial.variable]; !ok {
		return &ParseError{
			llmMessage: fmt.Sprintf("не объявлен аргумент %s", monomial.variable),
			message:    "undefined arg",
		}
	}
	return nil
}

func checkInterpretation(interpret Interpretation, constructorArity map[string]int) *ParseError {
	// TODO: check for duplicate args

	arity, ok := constructorArity[interpret.name]
	if !ok {
		return &ParseError{
			llmMessage: "конструктор отсутствует в правилах trs",
			message:    "excess interpretation",
		}
	}

	if arity != len(interpret.args) {
		return &ParseError{
			llmMessage: fmt.Sprintf("ожидался конструктор арности %d, получен арности %d", arity, len(interpret.args)),
			message:    "wrong func interpretation arity",
		}
	}

	return nil
}

func toMap(slice []string) map[string]struct{} {
	res := make(map[string]struct{})
	for _, el := range slice {
		res[el] = struct{}{}
	}
	return res
}
