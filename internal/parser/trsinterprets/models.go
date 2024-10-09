package trsinterprets

type Interpretation struct {
	Name      string
	Args      []string
	Monomials []Monomial
}

// NOTE: one of {constant, factors}
type Monomial struct {
	Constant *int
	Factors  *[]Factor
}

func NewConstantMonomial(v int) Monomial {
	return Monomial{
		Constant: &v,
		Factors:  nil,
	}
}

func NewProductMonomial(factors []Factor) Monomial {
	return Monomial{
		Constant: nil,
		Factors:  &factors,
	}
}

type Factor struct {
	Variable    string
	Coefficient int
	Power       int
}
