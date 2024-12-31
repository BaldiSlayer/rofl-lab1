import yaml
import faiss
import numpy as np

from sentence_transformers import SentenceTransformer

import app.core.vector_db.questions_preprocessing as question_preprocessor
import app.config.config as config
import app.schemas.questions as schemas


def convex_indexes(q_idx: int, counts: list[int]) -> (int, int):
    # индекс элемента базы знаний (у нас несколько вопросов на один ответ)
    elem_index = 0

    for item in counts:
        if q_idx < item:
            # возвращаем элемент базы знаний в котором содержится вопрос с номером q_idx и
            # номер вопроса внутри этого элемента
            return elem_index, q_idx

        q_idx -= item

        elem_index += 1

    raise ValueError("q_idx is out of bounds")


class FaissDB:
    model = None
    index = None

    MAX_CONTEXT_SIZE: int

    def __init__(self, sentence_transformer_name: str):
        self.model = SentenceTransformer(sentence_transformer_name)

        with open('data.yaml', 'r') as file:
            self.data = yaml.safe_load(file)

        self.elem_index_questions = []
        for item in self.data:
            self.elem_index_questions.append(len(item['questions']))

        self.index = faiss.read_index('vectorized_data.faiss')

        self.MAX_CONTEXT_SIZE = int(0.025 * len(self.data))

    def get_knowledge_base_elem(self, ans_pos: int, question_pos: int):
        elem = self.data[ans_pos]

        return schemas.QuestionAnswer(question=elem["questions"][question_pos], answer=elem["answer"])

    def get_context(self, distances, indices, similarity_threshold, k_max):
        context = []
        kb_items_idxes_set = set()

        for i in range(0, k_max):
            if distances[0][i] < similarity_threshold:
                break

            ans_pos, question_pos = convex_indexes(indices[0][i], self.elem_index_questions)

            # такой ответ еще не был добавлен в контекст
            if ans_pos not in kb_items_idxes_set:
                kb_items_idxes_set.add(ans_pos)

                context.append(self.get_knowledge_base_elem(ans_pos, question_pos))

        # TODO to not to dict
        return context

    def search_similar(self, query, similarity_threshold=0.5):
        """
        Dynamic search for similar objects based on similarity threshold.
        :param query: query string
        :param k_max: maximum number of results to return
        :param similarity_threshold: threshold for similarity to dynamically adjust k
        :return: list of similar objects
        """

        query_embedding = self.model.encode([question_preprocessor.prepocess_question(query)])

        faiss.normalize_L2(query_embedding)

        distances, indices = self.index.search(np.array(query_embedding), 2*self.MAX_CONTEXT_SIZE)

        return self.get_context(distances, indices, similarity_threshold, self.MAX_CONTEXT_SIZE)


faiss_db = FaissDB(
    config.SingletonConfig.get_instance().get_sentence_transformer_name(),
)


def process_question(question: str) -> list[schemas.QuestionAnswer]:
    """
    Process the question, finding similar knowledge base elements
    :param question:
    :return: nearest objects
    """

    similar_objects = faiss_db.search_similar(
        question,
    )

    return similar_objects
