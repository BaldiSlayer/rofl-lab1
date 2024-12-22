import faiss
import yaml
import sys
import string

from nltk import download
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize

from sentence_transformers import SentenceTransformer

import app.core.vector_db.text_translator as text_translator
import app.config.config as config


download('stopwords')
download('punkt_tab')

stop_words = set(stopwords.words('russian'))


def prepocess_question(lang_translator, question: str) -> str:
    question = question.strip()

    # Убираем пунктуацию и переводим текст в нижний регистр
    translator = str.maketrans('', '', string.punctuation)
    text = question.translate(translator).lower()

    words = word_tokenize(text)

    # Удаляем стоп слова
    filtered_words = [word for word in words if word not in stop_words]

    return lang_translator.translate_text(' '.join(filtered_words))


def create_embeddings(model, data):
    translator = text_translator.translator

    texts = []

    for item in data:
        for question in item["questions"]:
            texts.append(prepocess_question(translator, question))

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
