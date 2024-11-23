from typing import List
import yaml
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
    def should_include(contexts: List[dict]):
        contexts_questions = [i["question"] for i in contexts]
        contexts_set = {c for c in contexts_questions}

        nonlocal inclusion_list

        for question in inclusion_list:
            if question not in contexts_set:
                raise Exception(f"There is no question \"{question}\" in context {contexts_questions}")

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
            raise Exception(f"Error in test \"{self.question}\": {str(e)}")


def create_test(test_data: dict) -> Test:
    if 'should_include' in test_data:
        return Test(
            test_data["question"],
            should_include_wrapper(test_data["should_include"]),
        )

    raise Exception("There is no such test type")


def create_tests() -> List[Test]:
    with open("tests.yaml", 'r') as file:
        tests_data = yaml.safe_load(file)

    return [create_test(test_data) for test_data in tests_data]


def main():
    try:
        tests = create_tests()

        for i, test in enumerate(tests):
            test.check()
            print(f"Test {i + 1}/{len(tests)} passed")

        print("All tests passed successfully")
    except Exception as e:
        raise e


if __name__ == "__main__":
    main()
