from sentence_transformers import SentenceTransformer
import faiss
import numpy as np


def initialize_model():
    return SentenceTransformer('sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2')


def create_embeddings(model, data):
    texts = [item["question"] + " " + item["answer"] for item in data]
    return model.encode(texts)


def create_faiss_index(embeddings):
    dimension = embeddings.shape[1]
    index = faiss.IndexFlatL2(dimension)
    index.add(embeddings)
    return index


def search_similar(model, index, query, data, k=2):
    query_embedding = model.encode([query])
    D, I = index.search(np.array(query_embedding), k)
    return [data[idx] for idx in I[0]]


def process_questions(questions_list):
    model = initialize_model()
    embeddings = create_embeddings(model, questions_list)
    index = create_faiss_index(embeddings)

    results = []
    for item in questions_list:
        query = item["question"]
        similar_objects = search_similar(model, index, query, questions_list)
        result_str = f"Запрос: {query}\nПохожие объекты:\n"
        for obj in similar_objects:
            result_str += f"Вопрос: {obj['question']}, Ответ: {obj['answer']}\n"
        results.append(result_str)

    return results


# Пример использования
questions = [
    {"question": "Что такое ТФЯ?", "answer": "теория формальных языков"},
    {"question": "Что такое НЛП?", "answer": "обработка естественного языка"},
    {"question": "На какие языки перевод?", "answer": "перевод на русский язык и на английский языки"},
]

results = process_questions(questions)

for result in results:
    print(result)
    print("-" * 50)