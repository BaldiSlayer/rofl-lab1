from typing import List
import yaml
import openapi_client
from openapi_client.models.search_similar_request import SearchSimilarRequest

configuration = openapi_client.Configuration(
    host = "http://localhost:8100"
)

def get_similar(question: str) -> List[dict]:
    """
    Sends a request to the LM module API to find the nearest knowledge base elements for a given question

    :param question: The question for which the nearest ones will be searched in the knowledge base
    :return:
    """
    with openapi_client.ApiClient(configuration) as api_client:
        api_instance = openapi_client.QuestionsApi(api_client)

        try:
            api_response = api_instance.api_process_questions_process_questions_post(
                SearchSimilarRequest.from_dict({"question": question})
            )
            return api_response.to_dict()['result']
        except Exception as e:
            print("Exception when calling QuestionsApi->api_process_questions_process_questions_post: %s\n" % e)


def question_preprocessing(question: str) -> str:
    return question.strip()


def should_include_wrapper(inclusion_list: List[str]):
    def should_include(contexts: List[dict]):
        contexts_questions = [question_preprocessing(i["question"]) for i in contexts]
        contexts_set = {c for c in contexts_questions}

        nonlocal inclusion_list

        for question in inclusion_list:
            assert question.strip() in contexts_set, f"There is no question \"{question}\" in context {contexts_questions}"

        return True

    return should_include


class Test:
    def __init__(self, question: str, checker: callable):
        """
        Инициализация класса Test.

        :param question: Question from a testcase
        :param checker: Function that verifies compliance with the conditions
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
    tests = create_tests()

    for i, test in enumerate(tests):
        test.check()
        print(f"Test {i + 1}/{len(tests)} passed")

    print("All tests passed successfully")



if __name__ == "__main__":
    main()
