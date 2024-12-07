import re
from typing import List, Tuple

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
    res = question.strip().replace("\n", ' ')

    return re.sub(r'\s+', "", res)


def should_include_checker(question: str, should_include: List[str]):
    """
    should_include_checker sends request to get_similar route and checks if all items of
    should_include list are in context for question
    :param question: emulated user question
    :param should_include: list of questions that are necessary to be in context
    :return:
    """
    contexts = get_similar(question)

    contexts_answers = [question_preprocessing(i["answer"]) for i in contexts]
    contexts_set = {c for c in contexts_answers}

    for value in should_include:
        processed = question_preprocessing(value)

        assert question_preprocessing(value) in contexts_set, f"There is no question \"{processed}\" in context " \
                                                              f"answers {contexts_answers}"


def should_be_in_percentile_checker(question: str, should_be_in_percentile: List[Tuple[str, float]]):
    """
    should_be_in_percentile_checker sends request to get_similar route and checks if all
    items of should_be_in_percentile are in theirs percentiles.
    :param question: emulated user question
    :param should_be_in_percentile: list of (question, percentile)
    :return:
    """
    contexts = get_similar(question)

    contexts_questions = [question_preprocessing(i["question"]) for i in contexts]

    for value in should_be_in_percentile:
        percentile = contexts_questions.index(value[0])

        assert percentile != -1, f"There is no question \"{value}\" in context {contexts_questions}"

        assert value[1] >= (percentile / len(contexts_questions)), f" Question \"{value}\" is not in {value[1]} " \
                                                                   f"percentile in context {contexts_questions} "
