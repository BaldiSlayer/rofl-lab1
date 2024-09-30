from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional

from DB.faissDB import process_questions
from Mistral.mistral import get_chat_response
#uvicorn main:app --reload
app = FastAPI(title="TFL")

# Модель запроса для process_questions
class ProcessQuestionsRequest(BaseModel):
    questions_list: List[str]
    use_saved: Optional[bool] = False
    filename: Optional[str] = "vectorized_data"

# Маршрут для process_questions
@app.post("/process_questions")
def api_process_questions(request: ProcessQuestionsRequest):
    result = process_questions(
        questions_list=request.questions_list,
        use_saved=request.use_saved,
        filename=request.filename
    )
    return {"result": result}

# Модель запроса для get_chat_response
class GetChatResponseRequest(BaseModel):
    prompt: str
    context: Optional[str] = None
    model: Optional[str] = "open-mistral-7b"

# Маршрут для get_chat_response
@app.post("/get_chat_response")
def api_get_chat_response(request: GetChatResponseRequest):
    response = get_chat_response(
        prompt=request.prompt,
        context=request.context,
        model=request.model
    )
    return {"response": response}