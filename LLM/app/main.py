from fastapi import FastAPI

from app.routers.questions import router
import app.config.config as config
import app.core.vector_db.faiss_db as faiss_db
from app.core.readiness_probe.readiness_probe import ReadinessProbe


config.SingletonConfig.get_instance()

# проинициализировали конфиг и ставим ready в True
# так как загрузили индекс в память и готовы к обработке запросов
ReadinessProbe.set_value(True)

# uvicorn main:app --reload
app = FastAPI(
    title="LLM Proxy",
    version="0.0.1",
)

app.include_router(router)
