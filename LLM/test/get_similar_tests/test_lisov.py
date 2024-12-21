from test.get_similar_tests.common import should_include_checker, should_be_in_percentile_checker


def test_nearest_nka():
    question = "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?"

    should_include = [
        0,
    ]

    should_include_checker(question, should_include)


def test_nearest_shortest_word():
    question = "Дан регулярный язык, опиши алгоритм нахождения кратчайшего слова, принадлежащего этому регулярному " \
               "языку "

    should_include = [
        1,
    ]

    should_include_checker(question, should_include)


def test_nearest_shortest_word_2():
    question = "Алгоритм нахождения кратчайшего слова, принадлежащего этому регулярному языку"

    should_include = [
        1,
    ]

    should_include_checker(question, should_include)


def test_nearest_shortest_word_3():
    question = "А как мне найти кратчайшее слово, принадлежащее регулярному языку"

    should_include = [
        1,
    ]

    should_include_checker(question, should_include)
