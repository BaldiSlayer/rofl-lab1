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