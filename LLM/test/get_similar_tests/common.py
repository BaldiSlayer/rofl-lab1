from typing import List, Tuple

import numpy as np

import test.openapi_client as oc
from test.openapi_client.models.search_similar_request import SearchSimilarRequest

configuration = oc.Configuration(
    host="http://localhost:8100"
)


def get_similar(question: str) -> List[dict]:
    """
    Sends a request to the LM module API to find the nearest knowledge base elements for a given question

    :param question: The question for which the nearest ones will be searched in the knowledge base
    :return:
    """
    with oc.ApiClient(configuration) as api_client:
        api_instance = oc.QuestionsApi(api_client)

        try:
            api_response = api_instance.api_process_questions_process_questions_post(
                SearchSimilarRequest.from_dict({"question": question})
            )
            return api_response.to_dict()['result']
        except Exception as e:
            print("Exception when calling QuestionsApi->api_process_questions_process_questions_post: %s\n" % e)


def question_preprocessing(question: str) -> str:
    return question.strip()


def should_include_checker(question: str, should_include: List[str]):
    contexts = get_similar(question)

    contexts_questions = [question_preprocessing(i["question"]) for i in contexts]
    contexts_set = {c for c in contexts_questions}

    for question in should_include:
        assert question.strip() in contexts_set, f"There is no question \"{question}\" in context {contexts_questions}"


def should_be_in_percentile_checker(question: str, should_be_in_percentile: List[Tuple[str, int]]):
    return 1
