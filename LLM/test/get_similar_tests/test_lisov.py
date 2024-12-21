from test.get_similar_tests.common import should_include_checker, knowledge_base_data
import pytest


def test_nearest_nka(knowledge_base_data):
    question = "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?"

    should_include = [
        0,
    ]

    should_include_checker(knowledge_base_data, question, should_include)


def test_nearest_shortest_word(knowledge_base_data):
    question = "Дан регулярный язык, опиши алгоритм нахождения кратчайшего слова, принадлежащего этому регулярному " \
               "языку "

    should_include = [
        1,
    ]

    should_include_checker(knowledge_base_data, question, should_include)


def test_nearest_shortest_word_2(knowledge_base_data):
    question = "Алгоритм нахождения кратчайшего слова, принадлежащего этому регулярному языку"

    should_include = [
        1,
    ]

    should_include_checker(knowledge_base_data, question, should_include)


def test_nearest_shortest_word_3(knowledge_base_data):
    question = "А как мне найти кратчайшее слово, принадлежащее регулярному языку"

    should_include = [
        1,
    ]

    should_include_checker(knowledge_base_data, question, should_include)
