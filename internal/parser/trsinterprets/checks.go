package trsinterprets

import "fmt"

func checkMonomials(monomials []Monomial, args []string) *ParseError {
	c := monomialChecker{
		definedVars: toMap(args),
	}

	for _, monomial := range monomials {
		err := c.checkMonomial(monomial)
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
	if monomial.Factors == nil {
		return nil
	}
	if len(*monomial.Factors) == 0 {
		return &ParseError{
			llmMessage: "моном должен быть не пуст",
			message:    "empty monomial",
		}
	}

	for _, factor := range *monomial.Factors {
		err := c.checkFactor(factor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *monomialChecker) checkFactor(factor Factor) *ParseError {
	if _, ok := c.definedVars[factor.Variable]; !ok {
		return &ParseError{
			llmMessage: fmt.Sprintf("не объявлен аргумент %s", factor.Variable),
			message:    "undefined arg",
		}
	}
	return nil
}

func checkInterpretations(interprets []Interpretation, constructorArity map[string]int) *ParseError {
	c := interpretationsChecker{
		defined: map[string]struct{}{},
	}

	for _, interpretation := range interprets {
		err := c.checkInterpretation(interpretation, constructorArity)
		if err != nil {
			return err
		}
	}

	for expectedName := range constructorArity {
		if _, ok := c.defined[expectedName]; !ok {
			return &ParseError{
				llmMessage: fmt.Sprintf("не хватает интерпретации для конструктора %s", expectedName),
				message:    "no sufficient interpretation",
			}
		}
	}

	return nil
}

type interpretationsChecker struct {
	defined map[string]struct{}
}

type interpretationChecker = func(Interpretation, map[string]int) *ParseError

func (c *interpretationsChecker) checkInterpretation(interpret Interpretation,
	constructorArity map[string]int) *ParseError {

	checkers := []interpretationChecker{
		c.checkDuplicateInterpretation,
		c.checkExcessInterpretation,
		c.checkInterpretationArity,
		c.checkDuplicateArgumentName,
	}

	for _, checker := range checkers {
		err := checker(interpret, constructorArity)
		if err != nil {
			return err
		}
	}

	return nil

}

func (c *interpretationsChecker) checkDuplicateInterpretation(interpret Interpretation, _ map[string]int) *ParseError {
	if _, ok := c.defined[interpret.Name]; ok {
		return &ParseError{
			llmMessage: fmt.Sprintf("интерпретация конструктора %s задана повторно", interpret.Name),
			message:    "duplicate interpretation",
		}
	}
	c.defined[interpret.Name] = struct{}{}
	return nil
}

func (c *interpretationsChecker) checkExcessInterpretation(interpret Interpretation,
	constructorArity map[string]int) *ParseError {

	_, ok := constructorArity[interpret.Name]
	if !ok {
		return &ParseError{
			llmMessage: fmt.Sprintf("конструктор %s отсутствует в правилах trs", interpret.Name),
			message:    "excess interpretation",
		}
	}

	return nil
}

func (c *interpretationsChecker) checkInterpretationArity(interpret Interpretation,
	constructorArity map[string]int) *ParseError {

	expectedArity, _ := constructorArity[interpret.Name]
	if expectedArity != len(interpret.Args) {
		return &ParseError{
			llmMessage: fmt.Sprintf("неверная арность конструктора %s: "+
				"ожидалась арность %d, получена арность %d", interpret.Name, expectedArity, len(interpret.Args)),
			message: "wrong func interpretation arity",
		}
	}
	return nil
}

func (c *interpretationsChecker) checkDuplicateArgumentName(interpret Interpretation, _ map[string]int) *ParseError {
	args := map[string]struct{}{}
	for _, arg := range interpret.Args {
		if _, ok := args[arg]; ok {
			return &ParseError{
				llmMessage: fmt.Sprintf(
					"в интерпретации конструктора %s повторно объявлена переменная %s",
					interpret.Name,
					arg,
				),
				message: "duplicate argument name",
			}
		}
		args[arg] = struct{}{}
	}
	return nil
}

func toMap(slice []string) map[string]struct{} {
	res := make(map[string]struct{}, len(slice))
	for _, el := range slice {
		res[el] = struct{}{}
	}
	return res
}
