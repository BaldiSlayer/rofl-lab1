import yaml
import faiss
from sentence_transformers import SentenceTransformer
import numpy as np
import pickle
import os

import app.core.vector_db.text_translator as text_translator
import app.config.config as config
import app.schemas.questions as schemas

faiss_db = None


class FaissDB:
    model = None
    index = None

    def __init__(self, sentence_transformer_name: str):
        self.model = SentenceTransformer(sentence_transformer_name)

        # база знаний
        with open('data.yaml', 'r') as file:
            self.data = yaml.safe_load(file)

        self.index = faiss.read_index('vectorized_data.faiss')

    def search_similar(self, query, k_max=2, similarity_threshold=0.69):
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

        D, I = self.index.search(np.array(query_embedding), k_max)

        dynamic_k = 1
        for i in range(1, k_max):
            if D[0][i] < similarity_threshold:
                break
            dynamic_k += 1

        return [self.data[idx] for idx in I[0][:dynamic_k]]

    @staticmethod
    def create_embeddings(model, data):
        translator = text_translator.translator
        texts = [translator.translate_text(item["question"] + " " + item["answer"]) for item in data]
        embeddings = model.encode(texts)
        faiss.normalize_L2(embeddings)
        return embeddings

    @staticmethod
    def create_faiss_index(embeddings):
        dimension = embeddings.shape[1]
        index = faiss.IndexFlatIP(dimension)
        faiss.normalize_L2(embeddings)
        index.add(embeddings)
        return index

    @staticmethod
    def save_vectorized_data(data, embeddings, index, filename):
        with open(f"{filename}_data.pkl", "wb") as f:
            pickle.dump(data, f)
        np.save(f"{filename}_embeddings.npy", embeddings)
        faiss.write_index(index, f"{filename}_index.faiss")

    @staticmethod
    def load_vectorized_data(filename):
        with open(f"{filename}_data.pkl", "rb") as f:
            data = pickle.load(f)
        embeddings = np.load(f"{filename}_embeddings.npy")
        index = faiss.read_index(f"{filename}_index.faiss")
        return data, embeddings, index

    def process_questions(self, questions_list, use_saved=False, filename="vectorized_data"):
        model = self.model
        if use_saved and os.path.exists(f"{filename}_data.pkl"):
            data, embeddings, index = self.load_vectorized_data(filename)
        else:
            data = questions_list
            embeddings = self.create_embeddings(model, data)
            index = self.create_faiss_index(embeddings)
            self.save_vectorized_data(data, embeddings, index, filename)

        results = []
        for item in questions_list:
            query = item["question"]
            similar_objects = self.search_similar(model, index, query, data)
            result_str = f"Запрос: {query}\nПохожие объекты:\n"
            for obj in similar_objects:
                result_str += f"Вопрос: {obj['question']}, Ответ: {obj['answer']}\n"
            results.append(result_str)

        return results

    def add_new_questions(self, new_questions, filename="vectorized_data"):
        model = self.model
        if os.path.exists(f"{filename}_data.pkl"):
            data, embeddings, index = self.load_vectorized_data(filename)
        else:
            raise FileNotFoundError(f"No existing data found at {filename}. Please initialize the database first.")

        new_embeddings = self.create_embeddings(model, new_questions)
        updated_data = data + new_questions
        updated_embeddings = np.vstack((embeddings, new_embeddings))
        index.add(new_embeddings)
        self.save_vectorized_data(updated_data, updated_embeddings, index, filename)
        print(f"Added {len(new_questions)} new questions to the database.")


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

    conf = config.SingletonConfig.get_instance()
    similar_objects = faiss_db.search_similar(question)
    return [make_question_answer(qa_item) for qa_item in similar_objects]
