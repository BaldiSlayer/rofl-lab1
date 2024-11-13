import faiss
import yaml
import sys

import app.core.vector_db.faiss_db as faiss_db


def create_embeddings(model, data, output_filename) -> None:
    embeddings = faiss_db.create_embeddings(model, data)

    index = faiss_db.create_faiss_index(embeddings)

    faiss.write_index(index, output_filename)


if __name__ == "__main__":
    with open(sys.argv[1], 'r', encoding='utf-8') as file:
        content = yaml.safe_load(file)

    for item in content:
        if 'question' in item:
            item['question'] = item['question'].replace('\n', ' ')
        if 'answer' in item:
            item['answer'] = item['answer'].replace('\n', ' ')

    create_embeddings(
        faiss_db.SingletonSentenceTransformer.get_instance(),
        content,
        sys.argv[2],
    )
