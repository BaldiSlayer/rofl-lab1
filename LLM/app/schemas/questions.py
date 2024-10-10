from typing import List, Optional
from pydantic import BaseModel

# Модель запроса для process_questions
class ProcessQuestionsRequest(BaseModel):
    questions_list: List[str]
    use_saved: Optional[bool] = False
    filename: Optional[str] = "vectorized_data"


# Модель запроса для get_chat_response
class GetChatResponseRequest(BaseModel):
    prompt: str
    context: Optional[str] = None
    model: Optional[str] = "open-mistral-7b"


# Модель для вопроса и ответа
class QuestionAnswer(BaseModel):
    question: str
    answer: str


# Модель запроса для добавления новых вопросов
class AddQuestionsRequest(BaseModel):
    new_questions: List[QuestionAnswer]
    filename: Optional[str] = "vectorized_data"