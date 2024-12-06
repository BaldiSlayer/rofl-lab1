from test.get_similar_tests.common import should_include_checker


def test_nearest_example_1():
    question = "Что такое счетчиковая машина?"
    should_include = [
        "Что такое счетчиковая машина?"
    ]

    should_include_checker(question, should_include)
