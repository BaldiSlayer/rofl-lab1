from typing import List, Optional
from pydantic import BaseModel


# Модель для вопроса и ответа
class QuestionAnswer(BaseModel):
    question: str
    answer: str


# Модель запроса для process_questions
class ProcessQuestionsRequest(BaseModel):
    questions_list: List[QuestionAnswer]
    use_saved: Optional[bool] = False
    filename: Optional[str] = "vectorized_data"


# Модель запроса для get_chat_response
class GetChatResponseRequest(BaseModel):
    prompt: str
    context: Optional[str] = None
    model: Optional[str] = "open-mistral-7b"


# Модель запроса для добавления новых вопросов
class AddQuestionsRequest(BaseModel):
    new_questions: List[QuestionAnswer]
    filename: Optional[str] = "vectorized_data"


# Модель запроса для сохранения векторных данных
class SaveVectorizedDataRequest(BaseModel):
    data: List[QuestionAnswer]  # Список вопросов и ответов
    embeddings: List[List[float]]  # Эмбеддинги в виде списка списков
    filename: Optional[str] = "vectorized_data"  # Имя файла для сохранения
