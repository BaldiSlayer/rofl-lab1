package trsparser

import (
	"testing"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
)

func TestPeanoDto(t *testing.T) {
	f := make([]trsinterprets.Factor, 4)
	f[0] = trsinterprets.Factor{
		Coefficient: 1,
		Power:       1,
		Variable:    "a",
	}
	f[1] = trsinterprets.Factor{
		Coefficient: 2,
		Power:       1,
		Variable:    "b",
	}
	f[2] = trsinterprets.Factor{
		Coefficient: 1,
		Power:       2,
		Variable:    "c",
	}
	f[3] = trsinterprets.Factor{
		Coefficient: 3,
		Power:       4,
		Variable:    "d",
	}

	m := make([]trsinterprets.Monomial, 2)
	two := 2
	m[0] = trsinterprets.Monomial{
		Constant: &two,
	}
	m[1] = trsinterprets.Monomial{
		Factors: &f,
	}

	result := toInterpretationDTO(m)
	if result != "2+a*2*b*c**2*3*d**4" {
		t.Errorf("not expected: %s", result)
	}
}
