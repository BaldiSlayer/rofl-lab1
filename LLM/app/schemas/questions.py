from typing import List, Mapping, Optional
from pydantic import BaseModel


class QuestionAnswer(BaseModel):
    question: str
    answer: str


class SearchSimilarRequest(BaseModel):
    question: str

class SearchSimilarResponse(BaseModel):
    result: List[QuestionAnswer]


# Модель запроса для get_chat_response
class GetChatResponseRequest(BaseModel):
    prompt: str
    context: Optional[str] = None
    model: Optional[str] = "mistral-large-2411"


# Модель запроса для formalize
class FormalizeRequest(BaseModel):
    question: str
