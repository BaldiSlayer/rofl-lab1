from sentence_transformers import SentenceTransformer
import faiss
import numpy as np
import pickle
import os
from googletrans import Translator

translator = Translator()

def initialize_model():
    return SentenceTransformer('sentence-transformers/all-mpnet-base-v2')


def translate_text(text, dest_lang='en'):
    return translator.translate(text, dest=dest_lang).text


def create_embeddings(model, data):
    texts = [translate_text(item["question"] + " " + item["answer"]) for item in data]
    return model.encode(texts)


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
    # get the query vector
    query_embedding = model.encode([translate_text(query)])

    # perform a search with the maximum value of k
    D, I = index.search(np.array(query_embedding), k_max)

    # Dynamically determine k depending on the distances
    dynamic_k = 1  # At least one result is always returned
    for i in range(1, k_max):
        if D[0][i] - D[0][i - 1] > similarity_threshold:
            break
        dynamic_k += 1

    # Return only those objects that satisfy the dynamic k
    return [data[idx] for idx in I[0][:dynamic_k]]


def save_vectorized_data(data, embeddings, index, filename):
    """
    Save vectorized data, embeddings, and FAISS index
    :param data: original data
    :param embeddings: embeddings of the data
    :param index: FAISS index
    :param filename: base filename to save the data
    """
    with open(f"{filename}_data.pkl", "wb") as f:
        pickle.dump(data, f)

    np.save(f"{filename}_embeddings.npy", embeddings)

    faiss.write_index(index, f"{filename}_index.faiss")


def load_vectorized_data(filename):
    """
    Load vectorized data, embeddings, and FAISS index
    :param filename: base filename to load the data
    :return: data, embeddings, index
    """
    with open(f"{filename}_data.pkl", "rb") as f:
        data = pickle.load(f)

    embeddings = np.load(f"{filename}_embeddings.npy")

    index = faiss.read_index(f"{filename}_index.faiss")

    return data, embeddings, index


def process_questions(questions_list, use_saved=False, filename="vectorized_data"):
    """
    Process the questions, optionally using saved vectorized data
    :param questions_list: questions
    :param use_saved: whether to use saved vectorized data
    :param filename: base filename for saved data
    :return: nearest objects
    """
    model = initialize_model()

    if use_saved and os.path.exists(f"{filename}_data.pkl"):
        data, embeddings, index = load_vectorized_data(filename)
    else:
        data = questions_list
        embeddings = create_embeddings(model, data)
        index = create_faiss_index(embeddings)
        save_vectorized_data(data, embeddings, index, filename)

    results = []
    for item in questions_list:
        query = item["question"]
        similar_objects = search_similar(model, index, query, data)
        result_str = f"Запрос: {query}\nПохожие объекты:\n"
        for obj in similar_objects:
            result_str += f"Вопрос: {obj['question']}, Ответ: {obj['answer']}\n"
        results.append(result_str)

    return results


def add_new_questions(new_questions, filename="vectorized_data"):
    """
    Add new questions to the existing vectorized database and update the embeddings and FAISS index.
    :param new_questions: list of new questions in the format [{"question": "text", "answer": "text"}]
    :param filename: base filename for saved data
    """
    model = initialize_model()

    # Load existing data
    if os.path.exists(f"{filename}_data.pkl"):
        data, embeddings, index = load_vectorized_data(filename)
    else:
        raise FileNotFoundError(f"No existing data found at {filename}. Please initialize the database first.")

    # Create embeddings for new questions
    new_embeddings = create_embeddings(model, new_questions)

    # Concatenate old and new data and embeddings
    updated_data = data + new_questions
    updated_embeddings = np.vstack((embeddings, new_embeddings))

    # Update FAISS index with new embeddings
    index.add(new_embeddings)

    # Save the updated data, embeddings, and FAISS index
    save_vectorized_data(updated_data, updated_embeddings, index, filename)

    print(f"Added {len(new_questions)} new questions to the database.")


