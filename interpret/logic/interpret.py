#!/usr/bin/env python3
import re
from sympy import simplify as sympy_simplify, expand, Symbol
from random import randint
from z3 import *
from typing import List

'''
  The simplify_expression function simplifies expressions. 
  Multiplication must be entered with an asterisk '*'.
  Exponentiation can be implemented as both x**2 and x^2. 
  However, it's important to note that 2^4*5 will be simplified to 80. 
  So, if an expression is in the exponent, it must be enclosed in parentheses '()'.
'''
def simplify_expression(expression):
    # Replace variables with SymPy symbols
    expression = re.sub(r'([a-z]+)', r'Symbol("\1")', expression)

    simplified_expression = expand(sympy_simplify(expression))
    return str(simplified_expression)


'''
  The to_prefix_notation function rewrites an expression in prefix notation. 
  It also creates a set of variables contained in the resulting expression.
'''
def to_prefix_notation(expression, index):
    variables_set = set()
    expression = expression.replace('**', '^')

    summands = expression.split(' + ')
    prefix_expression = '(+'
    for summand in summands:
        multipliers = summand.split('*')
        prefix_multipliers = '(*'
        for multiplier in multipliers:
            if multiplier.isnumeric():
                prefix_multipliers = prefix_multipliers + ' ' + multiplier
            elif multiplier.isalpha():
                variable = multiplier + index
                variables_set.add(variable)
                prefix_multipliers = prefix_multipliers + ' ' + variable
            else:
                power = multiplier.split('^')
                variable = power[0] + index
                variables_set.add(variable)
                prefix_power = ' (^ ' + variable + ' ' + power[1] + ')'
                prefix_multipliers += prefix_power
        prefix_expression = prefix_expression + ' ' + prefix_multipliers + ')'
    prefix_expression += ')'
    return prefix_expression, variables_set

def get_arguments(function):
    i = 2
    result = []
    start = i
    round_brackets = 0
    while i < len(function):
        if function[i] == '(':
            round_brackets += 1
        elif function[i] == ')':
            round_brackets -= 1
        elif function[i] == ',' and round_brackets == 0:
            result.append(function[start:i])
            start = i+1
        i += 1
    result.append(function[start: -1])
    return result

'''
     Replaces function arguments with new values.
     Returns the right-hand side of the expression.
'''
def change_arguments(function_definition, new_args):
    result = ""
    parts_of_definition = function_definition.split('=')
    old_args = get_arguments(parts_of_definition[0])
    for symbol in parts_of_definition[1]:
        if symbol.isalpha():
            result = result + '(' + new_args[old_args.index(symbol)] + ')'
        else:
            result += symbol
    return result

def is_constructor(c, constructors):
    for i in constructors:
        if c == i:
            return True
    return False

'''
    Finds the inner function that starts at the given 'start' index.
    Returns the index of the element after that function
'''
def get_internal_function(str, start, constants):
    if str[start] in constants:
        return str[start], start+1
    else:
        i = start + 2
        round_brackets = 1

        while round_brackets > 0:
            if str[i] == '(':
                round_brackets += 1
            elif str[i] == ')':
                round_brackets -= 1
            i += 1
        return str[start:i], i

def find_function_definition(function, grammar_rules):
    for rule in grammar_rules:
        if function[0] == rule[0]:
            return rule
def substitute_interpretation(function, constructors, grammar_rules, constants):
    function_definition = find_function_definition(function, grammar_rules)

    #Situations where the function is just a constant or a variable
    if function_definition == None:
        return function
    elif function in constants:
        return function_definition[2:]
    args = get_arguments(function)
    for i in range(len(args)):
        arg = args[i]
        j = 0
        new_arg = ""
        while j < len(arg):
            if is_constructor(arg[j], constructors):
                internal_function, j = get_internal_function(arg, j, constants)
                new_arg = new_arg + '(' + substitute_interpretation(internal_function, constructors, grammar_rules, constants) + ')'
            else:
                new_arg += arg[j]
                j+=1
        args[i] = simplify_expression(new_arg)
    return simplify_expression(change_arguments(function_definition, args))

def interpret(trs_variables: List[str], trs_rules: List[str], grammar_rules: List[str]) -> str:
    result = ""
    constants = set()
    constructors = []
    for rule in grammar_rules:
        constructors.append(rule[0])
        if rule[1] == '=':
            constants.add(rule[0])

    def random_line(k, n, terms, l):
        t = randint(0, k)
        if l == 1:
            s = terms[t][0] + '('
            for i in range(n[t]):
                s = s + trs_variables[0] + '0,'
            s = s[:-1] + ')'
            return s
        s = terms[t][0] + '('
        for i in range(n[t]):
            s = s + random_line(k, n, terms, l - 1) + ','
        s = s[:-1] + ')'
        return s

    def demo():
        k = 0
        n = []
        terms = []
        for s in grammar_rules:
            if s[1] == '(':
                k += 1
                n.append(s.count(',') + 1)
                terms.append(s)
        l = randint(1, 5)
        s1 = random_line(k - 1, n, terms, l)
        li = randint(1, 5)
        s2 = s1
        for i in range(li):
            for j in trs_rules:
                if s2.find(j[0]) != -1:
                    a = s2.find(j[0]) + 2
                    b = 0
                    for rule in range(k):
                        if terms[rule][0] == j[0]:
                            b = rule
                            break
                    var_n = n[b]
                    new_s = j[j.find('=') + 1:]
                    sub_s = ""
                    c = 1
                    t = 0
                    i = 0
                    while t == 0:
                        if s2[a] == ',':
                            if c == 1:
                                new_s = new_s.replace(terms[b][2 * i + 2], sub_s)
                                sub_s = ""
                                i += 1
                            else:
                                sub_s = sub_s + s2[a]
                        elif s2[a] == ')':
                            if c == 1:
                                new_s = new_s.replace(terms[b][2 * i + 2], sub_s)
                                sub_s = ""
                                t = 1
                            else:
                                sub_s = sub_s + s2[a]
                                c += 1
                        elif s2[a] == '(':
                            sub_s = sub_s + s2[a]
                            c -= 1
                        else:
                            sub_s = sub_s + s2[a]
                        a += 1
                    s2 = s2[:s2.find(j[0])] + new_s + s2[a:]
                    break
        s1 = s1.replace(trs_variables[0] + '0', trs_variables[0])
        s2 = s2.replace(trs_variables[0] + '0', trs_variables[0])
        return s1, s2

    # -----------------------------------------------------------------------------
    start_expressions = []
    end_expressions = []
    with open('lab1.txt', 'w') as f:
        for cr in range(len(trs_rules)):
            trs_rules[cr] = trs_rules[cr].replace(" ", "")

            start = trs_rules[cr][:(trs_rules[cr].find('='))]
            end = trs_rules[cr][(trs_rules[cr].find('=') + 1):]

            start = substitute_interpretation(start, constructors, grammar_rules, constants)
            end = substitute_interpretation(end, constructors, grammar_rules, constants)

            start_expression, start_variables_set = to_prefix_notation(start, str(cr))
            end_expression, end_variables_set = to_prefix_notation(end, str(cr))
            variables_set = start_variables_set | end_variables_set
            start_expressions.append(start_expression)
            end_expressions.append(end_expression)
            for v in variables_set:
                f.write("(declare-fun " + v + " () Int)\n")
            for v in variables_set:
                f.write("(assert (>= " + v + " 0))\n")
        for cr in range(len(trs_rules)):
            f.write("(assert (<= " + start_expressions[cr] + " " + end_expressions[cr] + "))\n")
        f.write("(check-sat)\n")
        f.write("(get-model)")

    with open('lab1.txt', 'r') as f:
        smt_code = f.read()
    solver = Solver()
    solver.add(parse_smt2_string(smt_code))
    solver_result = solver.check()

    if solver_result == sat:
        result += "Сounterexample:\n"
        model = solver.model()
        for decl in model:
            result += f"{decl} = {model[decl]}\n"
    elif solver_result == unsat:
        result += "Verification success\nDemo:\n"
        s1, s2 = demo()
        result += s1 + '\n' + s2 + '\n'
        s1 = substitute_interpretation(s1, constructors, grammar_rules, constants)
        s2 = substitute_interpretation(s2, constructors, grammar_rules, constants)
        result += s1 + '\n' + s2 + '\n'

        s1 = str(s1).replace(trs_variables[0], '1')
        s2 = str(s2).replace(trs_variables[0], '1')
        s1 = simplify_expression(s1)
        s2 = simplify_expression(s2)
        result += s1 + '\n' + s2 + '\n'
    else:
        return "Unknown"
    return result
