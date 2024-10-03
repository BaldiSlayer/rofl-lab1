package trsparser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOnEmptyInput(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse("")

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		"система должна содержать хотя бы одно правило переписывания и его интерпретацию",
		parseError.LlmMessage,
	)
}

func TestParsesBasicTrs(t *testing.T) {
	expectedRule := Rule{
		Lhs: Subexpression{
			Args: &[]interface{}{NewSubexpression(Subexpression{
				Args:   &[]interface{}{},
				Letter: Letter{
					IsVariable: true,
					Name:       "x",
				},
			})},
			Letter: Letter{
				IsVariable: false,
				Name:       "f",
			},
		},
		Rhs: Subexpression{
			Args: &[]interface{}{},
			Letter: Letter{
				IsVariable: true,
				Name:       "x",
			},
		},
	}
	expectedInterpretation := Interpretation{
		Args:      []string{"x"},
		Monomials: []Monomial{NewConstantMonomial(5)},
		Name:      "f",
	}

	trs, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f(x) = 5
`,
	)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []Interpretation{expectedInterpretation},
		Rules:           []Rule{expectedRule},
		Variables:       []string{"x"},
	}, *trs)
}

func TestParsesComplexTrs(t *testing.T) {
	trs, err := Parser{}.Parse(
		`variables = x, y
f(x, g(y)) = g(f(x, h))
f(y, x) = g(y)
-----
f(x, y) = 5*x{2} + 10 + y{120}
g(x) = xx{2}5*x
h = 123
`,
	)
	got, _ := json.Marshal(*trs)
	os.Stdout.Write(got)

	assert.NoError(t, err)
	assert.JSONEq(t, `{
   "interpretations" : [
      {
         "args" : [
            "x",
            "y"
         ],
         "monomials" : [
            {
               "factors" : [
                  {
                     "coefficient" : 5,
                     "power" : 2,
                     "variable" : "x"
                  }
               ]
            },
            {
               "constant" : 10
            },
            {
               "factors" : [
                  {
                     "power" : 120,
                     "variable" : "y"
                  }
               ]
            }
         ],
         "name" : "f"
      },
      {
         "args" : [
            "x"
         ],
         "monomials" : [
            {
               "factors" : [
                  {
                     "variable" : "x"
                  },
                  {
                     "power" : 2,
                     "variable" : "x"
                  },
                  {
                     "coefficient" : 5,
                     "variable" : "x"
                  }
               ]
            }
         ],
         "name" : "g"
      },
      {
         "args" : [],
         "monomials" : [
            {
               "constant" : 123
            }
         ],
         "name" : "h"
      }
   ],
   "rules" : [
      {
         "lhs" : {
            "args" : [
               {
                  "args" : [],
                  "letter" : {
                     "isVariable" : true,
                     "name" : "x"
                  }
               },
               {
                  "args" : [
                     {
                        "args" : [],
                        "letter" : {
                           "isVariable" : true,
                           "name" : "y"
                        }
                     }
                  ],
                  "letter" : {
                     "isVariable" : false,
                     "name" : "g"
                  }
               }
            ],
            "letter" : {
               "isVariable" : false,
               "name" : "f"
            }
         },
         "rhs" : {
            "args" : [
               {
                  "args" : [
                     {
                        "args" : [],
                        "letter" : {
                           "isVariable" : true,
                           "name" : "x"
                        }
                     },
                     {
                        "args" : [],
                        "letter" : {
                           "isVariable" : false,
                           "name" : "h"
                        }
                     }
                  ],
                  "letter" : {
                     "isVariable" : false,
                     "name" : "f"
                  }
               }
            ],
            "letter" : {
               "isVariable" : false,
               "name" : "g"
            }
         }
      },
      {
         "lhs" : {
            "args" : [
               {
                  "args" : [],
                  "letter" : {
                     "isVariable" : true,
                     "name" : "y"
                  }
               },
               {
                  "args" : [],
                  "letter" : {
                     "isVariable" : true,
                     "name" : "x"
                  }
               }
            ],
            "letter" : {
               "isVariable" : false,
               "name" : "f"
            }
         },
         "rhs" : {
            "args" : [
               {
                  "args" : [],
                  "letter" : {
                     "isVariable" : true,
                     "name" : "y"
                  }
               }
            ],
            "letter" : {
               "isVariable" : false,
               "name" : "g"
            }
         }
      }
   ],
   "variables" : [
      "x",
      "y"
   ]
}
`, string(got))
}

func newInt(v int) *int {
	return &v
}
