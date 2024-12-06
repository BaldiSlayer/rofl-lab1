from test.get_similar_tests.common import should_include_checker


def test_nearest_nka():
    question = "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?"
    should_include = [
        "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
    ]

    should_include_checker(question, should_include)


def test_nearest_nka_2():
    question = "Опиши язык, распознаваемый недетерминированным конечным автоматом (НКА)?"
    should_include = [
        "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
    ]

    should_include_checker(question, should_include)


def test_nearest_nka_3():
    question = "Опиши язык, распознаваемый НКА"
    should_include = [
        "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
    ]

    should_include_checker(question, should_include)


def test_nearest_nka_4():
    question = "Опиши язык, распознаваемый НКА"
    should_include = [
        "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
    ]

    should_include_checker(question, should_include)


def test_nearest_nka_5():
    question = "Язык НКА. Описание"
    should_include = [
        "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
    ]

    should_include_checker(question, should_include)


