from interpret import interpret


def test_good_1():
    grammar_rules = [
        "z=1",
        "s(x)=x+1",
        "f(x,y)=x+y*2",
        "g(x,y)=x*y*3"
    ]
    trs_rules = [
        "f(x,s(y))=s(f(x,y))",
        "f(x,z)=x",
        "g(x,s(y))=g(g(x,y),x)",
        "g(x,z)=z"
    ]
    trs_variables = [
        "y"
        "x"
    ]

    result = interpret(trs_variables, trs_rules, grammar_rules)

    print(result)

    assert result == \
           'Сounterexample:\nx0 = 1\nx1 = 1\nx2 = 1\nx3 = 1\ny0 = 1\ny2 = 1\n'

