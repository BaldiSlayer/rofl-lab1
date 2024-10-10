#!/usr/bin/env python3

from typing import List


def interpret(trs_variables: List[str], trs_rules: List[str], grammar_rules: List[str]) -> str:
    for el in trs_rules:
        print(el)
    for el in grammar_rules:
        print(el)
    for el in trs_variables:
        print(el)
    return "unsat"
