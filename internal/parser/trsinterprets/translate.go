package trsinterprets

import "strconv"

func TranslateInterpretation(m []Monomial) string{
	result := ""
	count := len(m)
	for i, monom := range m{
		if monom.Constant != nil{ // если это константа
			result += strconv.Itoa(*monom.Constant)
		}else{ // если моном
			countMonoms := len(*monom.Factors)
			for j, factor := range *monom.Factors{
				if factor.Coefficient != 1{
					result += strconv.Itoa(factor.Coefficient) + "*"
				}
				result += factor.Variable
				if factor.Power != 1{
					result += "**" + strconv.Itoa(factor.Power)
				}
				if j != countMonoms - 1{
					result += "*"
				}
			}
		}
		if i != count - 1{
			result += "+"
		}
	}
	return result
}