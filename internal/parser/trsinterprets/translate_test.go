package trsinterprets

import "testing"

func TestParserWithPeano(t *testing.T) {
	f := make([]Factor, 4)
	f[0] = Factor{
		Coefficient: 1,
		Power:       1,
		Variable:    "a",
	}
	f[1] = Factor{
		Coefficient: 2,
		Power:       1,
		Variable:    "b",
	}
	f[2] = Factor{
		Coefficient: 1,
		Power:       2,
		Variable:    "c",
	}
	f[3] = Factor{
		Coefficient: 3,
		Power:       4,
		Variable:    "d",
	}

	m := make([]Monomial, 2)
	two := 2
	m[0] = Monomial{
		Constant: &two,
	}
	m[1] = Monomial{
		Factors: &f,
	}

	result := TranslateInterpretation(m)
	if result != "2+a*2*b*c**2*3*d**4" {
		t.Errorf("not expected: %s", result)
	}
}
