from fastapi import FastAPI

from app.routers.questions import router
import app.config.config as config


# инициализируем конфиг
config.SingletonConfig.get_instance()

# uvicorn main:app --reload
app = FastAPI(
    title="LLM Proxy",
    version="0.0.1",
)

app.include_router(router)
