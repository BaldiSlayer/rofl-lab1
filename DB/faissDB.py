from sentence_transformers import SentenceTransformer
import faiss
import numpy as np

# Инициализация модели для получения векторных представлений
model = SentenceTransformer('sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2')

# Пример данных
data = [
    {"question": "Что такое ТФЯ?", "answer": "теория формальных языков"},
    {"question": "Что такое НЛП?", "answer": "обработка естественного языка"},
    {"question": "Как работает GPT?", "answer": "модель, основанная на трансформерах"},
]

# Создание списка текстов (соединяем вопрос и ответ для каждого объекта)
texts = [item["question"] + " " + item["answer"] for item in data]

# Генерация эмбеддингов для каждого объекта
embeddings = model.encode(texts)

# Преобразование вектора эмбеддингов в формат numpy
embeddings = np.array(embeddings)

# Создание индекса FAISS для поиска ближайших соседей
dimension = embeddings.shape[1]  # Размерность векторов
index = faiss.IndexFlatL2(dimension)  # Используем L2 расстояние (евклидово)
index.add(embeddings)  # Добавляем вектора в индекс

# Пример запроса
query = "что такое формальные языки?"

# Генерация эмбеддинга для запроса
query_embedding = model.encode([query])

# Поиск ближайших объектов в базе данных
D, I = index.search(np.array(query_embedding), k=2)  # Найти ближайшие 2 объекта

# Вывод результатов
print("Похожие объекты:")
for idx in I[0]:
    print(data[idx])