import faiss
import yaml
import sys

import app.core.vector_db.text_translator as text_translator
from sentence_transformers import SentenceTransformer
import app.config.config as config


def create_embeddings(model, data):
    translator = text_translator.translator

    texts = []

    for item in data:
        for question in item["questions"]:
            texts.append(translator.translate_text(question))

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


def init_embeddings(model, data, output_filename) -> None:
    """
    builds embeddings, saves them to index in and saves it to the file
    :param model:
    :param data:
    :param output_filename:
    :return:
    """
    embeddings = create_embeddings(model, data)

    index = create_faiss_index(embeddings)

    faiss.write_index(index, output_filename)


def main():
    cfg = config.SingletonConfig.get_instance()

    with open(sys.argv[1], 'r', encoding='utf-8') as file:
        content = yaml.safe_load(file)

    init_embeddings(
        SentenceTransformer(cfg.get_sentence_transformer_name()),
        content,
        sys.argv[2],
    )


if __name__ == "__main__":
    main()
