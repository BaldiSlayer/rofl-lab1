import faiss
from sentence_transformers import SentenceTransformer
import numpy as np

import app.core.vector_db.text_translator as text_translator
import app.config.config as config
import app.schemas.questions as schemas


class SingletonSentenceTransformer:
    _instance = None

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = SentenceTransformer('sentence-transformers/all-mpnet-base-v2')

        return cls._instance


def create_embeddings(model, data):
    translator = text_translator.SingletonTextTranslator.get_instance()

    # TODO
    texts = [translator.translate_text(item["question"] + " " + item["answer"]) for item in data]
    embeddings = model.encode(texts)

    faiss.normalize_L2(embeddings)

    return embeddings


def create_faiss_index(embeddings):
    dimension = embeddings.shape[1]
    index = faiss.IndexFlatIP(dimension)
    # Нормализуем векторы перед добавлением в индекс
    faiss.normalize_L2(embeddings)
    index.add(embeddings)

    return index


def search_similar(model, index, query, data, k_max=10, similarity_threshold=0.1):
    """
    Dynamic search for similar objects based on similarity threshold.
    :param model: sentence transformer model
    :param index: FAISS index
    :param query: query string
    :param data: original data (list of questions/answers)
    :param k_max: maximum number of results to return
    :param similarity_threshold: threshold for similarity to dynamically adjust k
    :return: list of similar objects
    """

    translator = text_translator.SingletonTextTranslator.get_instance()

    query_embedding = model.encode([translator.translate_text(query)])

    faiss.normalize_L2(query_embedding)

    # perform a search with the maximum value of k
    distances, indices = index.search(np.array(query_embedding), k_max)

    # Find the closest distance
    closest_distance = distances[0][0]

    # Dynamically determine k depending on the distances
    dynamic_k = 1  # At least one result is always returned
    for i in range(1, k_max):
        if distances[0][i] - closest_distance > similarity_threshold:
            break
        dynamic_k += 1

    # Return only those objects that satisfy the dynamic k
    return [data[idx] for idx in indices[0][:dynamic_k]]


def make_question_answer(qa_item: dict) -> schemas.QuestionAnswer:
    return schemas.QuestionAnswer(question=qa_item["question"], answer=qa_item["answer"])


def process_questions(question: str) -> list[schemas.QuestionAnswer]:
    """
    Process the questions, optionally using saved vectorized data
    :param question:
    :return: nearest objects
    """

    model = SingletonSentenceTransformer.get_instance()
    conf = config.SingletonConfig.get_instance()

    similar_objects = search_similar(
        model,
        conf.get_index(),
        question,
        conf.get_data(),
    )

    return [make_question_answer(qa_item) for qa_item in similar_objects]
