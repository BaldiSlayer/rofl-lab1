from fastapi import APIRouter
from ..utils.DB.faissDB import process_questions
from ..schemas.questions import ProcessQuestionsRequest, GetChatResponseRequest
from ..utils.Mistral.mistral import get_chat_response


router = APIRouter(
    tags=["Questions"]
)
# Маршрут для process_questions
@router.post("/process_questions")
def api_process_questions(request: ProcessQuestionsRequest):
    result = process_questions(
        questions_list=request.questions_list,
        use_saved=request.use_saved,
        filename=request.filename
    )
    return {"result": result}


# Маршрут для get_chat_response
@router.post("/get_chat_response")
def api_get_chat_response(request: GetChatResponseRequest):
    response = get_chat_response(
        prompt=request.prompt,
        context=request.context,
        model=request.model
    )
    return {"response": response}