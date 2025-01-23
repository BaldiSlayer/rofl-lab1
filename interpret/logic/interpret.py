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
                s = s + 'var,'
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
                pattern = j.split('=')[0].strip()
                replacement = j.split('=')[1].strip()
                s2 = replace_pattern(s2, pattern, replacement, trs_variables)
        s1 = s1.replace('var', trs_variables[0])
        s2 = s2.replace('var', trs_variables[0])
        return s1, s2

    def replace_pattern(s, pattern, replacement, variables):
        pattern_args = get_arguments(pattern)
        args_mapping = {}
        pattern_index = s.find(pattern.split('(')[0])
        while pattern_index != -1:
            if check_pattern_match(s, pattern_index, pattern, pattern_args, variables, args_mapping):
                new_replacement = replacement
                for arg in (variables):
                    if arg in args_mapping:
                        new_replacement = new_replacement.replace(arg, args_mapping[arg])


                end_index = find_end_index(s, pattern_index)
                s = s[:pattern_index] + new_replacement + s[end_index:]
                pattern_index = s.find(pattern.split('(')[0], pattern_index + len(new_replacement))
            else:
                pattern_index = s.find(pattern.split('(')[0], pattern_index + 1)

        return s

    def find_end_index(s, start_index):
        open_brackets = 0
        current_index = start_index
        while current_index < len(s):
            if s[current_index] == '(':
                open_brackets += 1
            elif s[current_index] == ')':
                open_brackets -= 1
                if open_brackets == 0:
                    return current_index + 1
            current_index += 1
        return current_index

    def check_pattern_match(s, index, pattern, pattern_args, variables, args_mapping):
        pattern_len = len(pattern.split('(')[0])
        if s[index:index + pattern_len] != pattern.split('(')[0]:
            return False

        args_start = index + pattern_len + 1
        for i, arg in enumerate(pattern_args):
            if arg in variables:
                arg_value = get_var_value(s, args_start, i)
                if arg in args_mapping and args_mapping[arg] != arg_value:
                    return False
                else:
                    args_mapping[arg]=arg_value
                    continue
            arg_value = get_var_value(s, args_start, i)
            if not check_pattern_match(arg_value, 0, arg, get_arguments(arg), variables, args_mapping):
                return False
        return True

    def get_var_value(s, start, arg_index):
        open_brackets = 0
        current_index = start
        while current_index < len(s):
            if s[current_index] == '(':
                open_brackets += 1
            elif s[current_index] == ')':
                open_brackets -= 1
            elif s[current_index] == ',' and open_brackets == 0:
                if arg_index == 0:
                    return s[start:current_index].strip()
                arg_index -= 1
                start=current_index+1
            current_index += 1
        s=s[start:current_index]
        start=0
        while open_brackets<0:
            s=s[start:s.rfind(")")]
            open_brackets += 1
        return s

    # -----------------------------------------------------------------------------
    start_expressions = []
    end_expressions = []
    starts = []
    ends = []
    with open('lab1.txt', 'w') as f:
        for cr in range(len(trs_rules)):
            trs_rules[cr] = trs_rules[cr].replace(" ", "")

            start = trs_rules[cr][:(trs_rules[cr].find('='))]
            end = trs_rules[cr][(trs_rules[cr].find('=') + 1):]

            start = substitute_interpretation(start, constructors, grammar_rules, constants)
            end = substitute_interpretation(end, constructors, grammar_rules, constants)

            start_expression, start_variables_set = to_prefix_notation(start, str(cr))
            end_expression, end_variables_set = to_prefix_notation(end, str(cr))
            for trs_var in trs_variables:
                if start.find(trs_var) != -1:
                    for var in start_variables_set:
                        if var[0]==trs_var[0]:
                            start = start.replace(trs_var, var)
                if end.find(trs_var) != -1:
                    for var in end_variables_set:
                        if var[0]==trs_var[0]:
                            end = end.replace(trs_var, var)

            mistake1=start_expression.find("^")
            while mistake1!=-1:
                mistake1-=1
                var= start_expression[mistake1+3:mistake1+5]
                mistake2=start_expression[mistake1+6:].find(')')+mistake1+6
                power=int(start_expression[mistake1+6:mistake2])
                powerstring='1'
                for i in range(power):
                    powerstring= '(* '+var+ ' ' + powerstring+')'
                start_expression=start_expression[:mistake1]+powerstring+start_expression[mistake2+1:]
                mistake1=start_expression.find("^")
            mistake1=end_expression.find("^")
            while mistake1!=-1:
                mistake1-=1
                var= end_expression[mistake1+3:mistake1+5]
                mistake2=end_expression[mistake1+6:].find(')')+mistake1+6
                power=int(end_expression[mistake1+6:mistake2])
                powerstring='1'
                for i in range(power):
                    powerstring= '(* '+var+ ' ' + powerstring+')'
                end_expression=end_expression[:mistake1]+powerstring+end_expression[mistake2+1:]
                mistake1=end_expression.find("^")
            #print(start, start_expression)
            #print(end, end_expression)
            variables_set = start_variables_set | end_variables_set
            start_expressions.append(start_expression)
            end_expressions.append(end_expression)
            starts.append(start)
            ends.append(end)
            for v in variables_set:
                f.write("(declare-fun " + v + " () Int)\n")
            for v in variables_set:
                f.write("(assert (> " + v + " 0))\n")
        assert_line = "(<= " + start_expressions[0] + " " + end_expressions[0] + ")"
        for cr in range(1, len(trs_rules)):
            assert_line = "(or " + assert_line + " (<= " + start_expressions[cr] + " " + end_expressions[cr] + "))"
        f.write("(assert " + assert_line + ")\n")
        f.write("(check-sat)\n")
        f.write("(get-model)")

    with open('lab1.txt', 'r') as f:
        smt_code = f.read()
    solver = Solver()
    solver.add(parse_smt2_string(smt_code))
    solver_result = solver.check()

    if solver_result == sat:
        result += "Verification is not successful\n"
        result += "Counterexample:\n"
        model = solver.model()
        n=len(starts)
        for i in range(n):
            changed_start = starts[i]
            changed_end = ends[i]
            counterexample_statement = "If "
            for declaration in model:
                decl=str(declaration)
                #print(declaration, decl, model[declaration])
                if changed_start.find(decl) != -1 or changed_end.find(decl) != -1:
                    changed_start = changed_start.replace(decl, str(model[declaration]))
                    changed_end = changed_end.replace(decl, str(model[declaration]))
                    counterexample_statement = counterexample_statement + decl + " = " + str(model[declaration]) + "\n"
            counterexample_statement = counterexample_statement + "Then " + trs_rules[i][:(trs_rules[i].find('='))] + " > " + trs_rules[i][(trs_rules[i].find('=') + 1):] +"\n"
            counterexample_statement = counterexample_statement + "Equals " + starts[i] + " > " + ends[i] +"\n"
            counterexample_statement = counterexample_statement + "Equals " + changed_start + " > " + changed_end +"\n"
            #print(changed_start, changed_end)
            counted_start = int(simplify_expression(changed_start))
            counted_end = int(simplify_expression(changed_end))
            if counted_start <= counted_end:
                counterexample_statement = counterexample_statement + "But, unfortunately, " + str(counted_start) + " <= " + str(counted_end) + "\n"
                result += counterexample_statement
    elif solver_result == unsat:
        result += "Verification success\nDemo:\n"
        s1, s2 = demo()
        while s1 == s2:
            s1, s2 = demo()
        result += 'Generated line: ' + s1 + '\nSimplified line: ' + s2 + '\n'
        s1 = substitute_interpretation(s1, constructors, grammar_rules, constants)
        s2 = substitute_interpretation(s2, constructors, grammar_rules, constants)
        result += 'Interpreted generated line: ' + s1 + '\nInterpreted simplified line: ' + s2 + '\n'
        s1 = str(s1).replace(trs_variables[0], '1')
        s2 = str(s2).replace(trs_variables[0], '1')
        s1 = simplify_expression(s1)
        s2 = simplify_expression(s2)
        result += 'Interpreted generated line weight: ' + s1 + '\nInterpreted simplified line weight: ' + s2 + '\n'
    else:
        return "Unknown"
    return result
