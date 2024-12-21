import re
import yaml
from typing import List, Tuple
import pytest

import test.openapi_client as oc
from test.openapi_client.models.search_similar_request import SearchSimilarRequest

configuration = oc.Configuration(
    host="http://localhost:8100"
)


@pytest.fixture(scope='session')
def knowledge_base_data() -> dict:
    with open('data.yaml', 'r', encoding='utf-8') as file:
        data = yaml.safe_load(file)

    return {i['id']: i for i in data}


def get_similar(question: str) -> List[dict]:
    """
    Sends a request to the LM module API to find the nearest knowledge base elements for a given question

    :param: question: The question for which the nearest ones will be searched in the knowledge base
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


def should_include_checker(knowledge_base_data, question: str, should_include: List[int]):
    """
    should_include_checker sends request to get_similar route and checks if all items of
    should_include list are in context for question
    :param knowledge_base_data:
    :param question: emulated user question
    :param should_include: list of answers ids that are necessary to be in context
    :return:
    """
    contexts = get_similar(question)

    contexts_answers = [i["answer"] for i in contexts]
    proceeded_contexts_answers = {question_preprocessing(i) for i in contexts_answers}

    for value_id in should_include:
        value = knowledge_base_data[value_id]['answer']
        processed = question_preprocessing(value)

        assert processed in proceeded_contexts_answers, \
            f"There is no answer \"{value}\" in context answers {contexts_answers}"


def should_be_in_percentile_checker(
    knowledge_base_data,
    question: str,
    should_be_in_percentile: List[Tuple[str, float]],
):
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
        ans_index = contexts_questions.index(value[0])

        assert ans_index != -1, f"There is no answer \"{value}\" in context {contexts_questions}"

        percentile = ans_index / len(contexts_questions)

        assert value[1] >= percentile, \
            f" Question \"{value}\" is not in {value[1]} percentile in context {contexts_questions}"
