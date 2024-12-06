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


def test_epsilon_closure_1():
    question = "Опиши алгоритм нахождения эпсилон замыкания для каждой из вершин автомата"
    should_include = [
        "Опиши алгоритм нахождения эпсилон замыкания для каждой из вершин автомата?",
    ]

    should_include_checker(question, should_include)


def test_epsilon_closure_2():
    question = "Как мне найти эпсилон замыкание для каждой из вершин автомата?"
    should_include = [
        "Опиши алгоритм нахождения эпсилон замыкания для каждой из вершин автомата?",
    ]

    should_include_checker(question, should_include)


def test_epsilon_closure_3():
    question = "Нахождение эпсилон замыкания для все вершин автомата?"
    should_include = [
        "Опиши алгоритм нахождения эпсилон замыкания для каждой из вершин автомата?",
    ]

    should_include_checker(question, should_include)


def test_prefix_languages_1():
    question = "Какой язык называется префиксным (беспрефиксным)"
    should_include = [
        "Какой язык называется префиксным (беспрефиксным)",
    ]

    should_include_checker(question, should_include)


def test_prefix_languages_2():
    question = "Дай определение префиксного языка"
    should_include = [
        "Какой язык называется префиксным (беспрефиксным)",
    ]

    should_include_checker(question, should_include)


def test_prefix_languages_3():
    question = "Дай определение беспрефиксного языка"
    should_include = [
        "Какой язык называется префиксным (беспрефиксным)",
    ]

    should_include_checker(question, should_include)


def test_prefix_languages_4():
    question = "Префиксный язык определение"
    should_include = [
        "Какой язык называется префиксным (беспрефиксным)",
    ]

    should_include_checker(question, should_include)


def test_prefix_languages_5():
    question = "Беспрефиксный язык определение"
    should_include = [
        "Какой язык называется префиксным (беспрефиксным)",
    ]

    should_include_checker(question, should_include)


def test_count_n_len_words_1():
    question = "Опиши алгоритм подсчета количества слов определенной длины в заданном регулярном языке"
    should_include = [
        "Опиши алгоритм подсчета количества слов определенной длины в заданном регулярном языке",
    ]

    should_include_checker(question, should_include)


def test_count_n_len_words_2():
    question = "Алгоритм подсчета количества слов определенной длины в заданном регулярном языке"
    should_include = [
        "Опиши алгоритм подсчета количества слов определенной длины в заданном регулярном языке",
    ]

    should_include_checker(question, should_include)


def test_count_n_len_words_3():
    question = "Как подсчитать количество слов определенной длины в заданном регулярном языке"
    should_include = [
        "Опиши алгоритм подсчета количества слов определенной длины в заданном регулярном языке",
    ]

    should_include_checker(question, should_include)

