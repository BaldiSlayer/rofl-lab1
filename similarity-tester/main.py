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


def should_include_wrapper(inclusion_list: List[str]):
    def should_include(contexts: List[str]):
        contexts_set = {c for c in contexts}

        nonlocal inclusion_list

        for question in inclusion_list:
            if question not in contexts_set:
                raise Exception(f"There is no question {question} in context {contexts}")

        return True

    return should_include



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

        try:
            return self._checker(answer)
        except Exception as e:
            raise e



def create_tests() -> List[Test]:
    return [
        Test("Регулярен ли язык Дика", should_include_wrapper(["что-то"]))
    ]


def main():
    tests = create_tests()

    for i, test in enumerate(tests):
        try:
            test.check()
        except Exception as e:
            raise e

        print(f"Пройден тест {i + 1}/{len(tests)}")



    print("Все тесты пройдены успешно")


if __name__ == "__main__":
    main()
