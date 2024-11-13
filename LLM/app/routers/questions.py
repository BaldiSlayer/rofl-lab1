from fastapi import APIRouter
from fastapi.responses import JSONResponse

import app.core.vector_db.faiss_db as faiss_db
import app.schemas.questions as schemas
import app.core.mistral.mistral as mistral_client
import app.core.formalize.formalize as formalize
from app.core.readiness_probe.readiness_probe import ReadinessProbe


router = APIRouter(
    tags=["Questions"]
)


@router.get("/ping")
def api_ping():
    """
    Ping ручка, нужна для healthcheck'а
    """

    if not ReadinessProbe.get_value():
        return JSONResponse(status_code=503, content={"status": "not ready"})

    return "OK"


@router.post("/search_similar", response_model=schemas.SearchSimilarResponse)
def api_search_similar(request: schemas.SearchSimilarRequest):
    """
    Ищет для вопроса похожие из базы знаний
    """

    result = faiss_db.process_questions(
        question=request.question,
    )

    return schemas.SearchSimilarResponse(result=result)


@router.post("/get_chat_response")
def api_get_chat_response(request: schemas.GetChatResponseRequest):
    """
    Ручка, которая ходит в Mistral API с заданным промптом, контекстом и моделью
    """

    # TODO add reponse model for openapi
    response = mistral_client.get_chat_response(
        prompt=request.prompt,
        context=request.context,
        model=request.model
    )

    return {"response": response}


@router.post("/formalize")
def api_formalize(request: schemas.FormalizeRequest):
    """
    Ручка для formalize
    """

    response = formalize.formalize(
        question=request.question
    )

    return {"response": response}
