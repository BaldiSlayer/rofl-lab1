import yaml
import faiss
from sentence_transformers import SentenceTransformer
import numpy as np

import app.core.vector_db.text_translator as text_translator
import app.config.config as config
import app.schemas.questions as schemas

faiss_db = None


def convex_indexes(q_idx: int, counts: list[int]) -> (int, int):
    elem_index = 0

    for item in counts:
        if q_idx < item:
            return elem_index, q_idx

        q_idx -= item

        elem_index += 1

    raise ValueError("q_idx is out of bounds")


class FaissDB:
    model = None
    index = None

    def __init__(self, sentence_transformer_name: str):
        self.model = SentenceTransformer(sentence_transformer_name)

        # база знаний
        with open('data.yaml', 'r') as file:
            self.data = yaml.safe_load(file)

        self.elem_index_questions = []
        for item in self.data:
            self.elem_index_questions.append(len(item['questions']))

        self.index = faiss.read_index('vectorized_data.faiss')

    def get_knowledge_base_elem(self, ans_pos: int, question_pos: int):
        print(ans_pos, question_pos)

        elem = self.data[ans_pos]

        return {"question": elem["questions"][question_pos], "answer": elem["answer"]}

    def get_context(self, distances, indices, similarity_threshold, k_max):
        context = []
        kb_items_idxes_set = set()

        closest_distance = distances[0][0]

        for i in range(0, k_max):
            if distances[0][i] - closest_distance > similarity_threshold:
                break

            ans_pos, question_pos = convex_indexes(indices[0][i], self.elem_index_questions)

            # такой вопрос еще не был добавлен в контекст
            if ans_pos not in kb_items_idxes_set:
                kb_items_idxes_set.add(ans_pos)

                context.append(self.get_knowledge_base_elem(ans_pos, question_pos))

        return context

    def search_similar(self, query, k_max=2, similarity_threshold=0.1):
        """
        Dynamic search for similar objects based on similarity threshold.
        :param query: query string
        :param k_max: maximum number of results to return
        :param similarity_threshold: threshold for similarity to dynamically adjust k
        :return: list of similar objects
        """

        translator = text_translator.translator

        query_embedding = self.model.encode([translator.translate_text(query)])

        faiss.normalize_L2(query_embedding)

        # perform a search with the maximum value of k
        distances, indices = self.index.search(np.array(query_embedding), k_max)

        return self.get_context(distances, indices, similarity_threshold, k_max)


def init_faiss_db():
    global faiss_db

    faiss_db = FaissDB(
        config.SingletonConfig.get_instance().get_sentence_transformer_name(),
    )


def make_question_answer(qa_item: dict) -> schemas.QuestionAnswer:
    return schemas.QuestionAnswer(question=qa_item["question"], answer=qa_item["answer"])


def process_question(question: str) -> list[schemas.QuestionAnswer]:
    """
    Process the question, finding similar knowledge base elements
    :param question:
    :return: nearest objects
    """

    similar_objects = faiss_db.search_similar(
        question,
    )

    return [make_question_answer(qa_item) for qa_item in similar_objects]
