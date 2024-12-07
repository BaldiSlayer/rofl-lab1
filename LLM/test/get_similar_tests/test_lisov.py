from test.get_similar_tests.common import should_include_checker, should_be_in_percentile_checker


def test_nearest_nka():
    question = "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?"
    should_include = [
        """Язык, распознаваемый недетерминированным конечным автоматом (НКА) – это
    все такие слова,  по которым существует хотя бы один путь из стартовой вершины
    в терминальную."""
    ]

    should_include_checker(question, should_include)
