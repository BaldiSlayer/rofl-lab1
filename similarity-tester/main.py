from typing import List
import openapi_client
from openapi_client.models.search_similar_request import SearchSimilarRequest

configuration = openapi_client.Configuration(
    host = "http://localhost:8100"
)

def get_similar(question: str):
    with openapi_client.ApiClient(configuration) as api_client:
        # Create an instance of the API class
        api_instance = openapi_client.QuestionsApi(api_client)

        try:
            # Api Process Questions
            api_response = api_instance.api_process_questions_process_questions_post(
                SearchSimilarRequest.from_dict({"question": question})
            )
            return api_response.to_dict()['result']
        except Exception as e:
            print("Exception when calling QuestionsApi->api_process_questions_process_questions_post: %s\n" % e)


class Test:
    def __init__(self, question: str, checker: callable):
        """
        Инициализация класса Test.

        :param question:
        :param checker: Функция, которая проверяет правильность ответа.
        """
        self.question = question
        self._checker = checker

    def check(self) -> bool:
        answer = get_similar(self.question)

        return self._checker(answer)


def h(a: dict):
    if a == {}:
        raise Exception("ожидался пустой словарь")

    return True


def create_tests() -> List[Test]:
    return [
        Test("Регулярен ли язык Дика", h)
    ]


def main():
    tests = create_tests()

    for i, test in enumerate(tests):
        try:
            test.check()
        except Exception as e:
            raise Exception(f'Тест с вопросом "{test.question}" не пройден: {str(e)}') from e

        print(f"Пройден тест {i + 1}/{len(tests)}")



    print("Все тесты пройдены успешно")


if __name__ == "__main__":
    main()
