FROM python:3.12-slim

WORKDIR /api

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

COPY /LLM/requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY /LLM .


CMD gunicorn main:app --workers 1 --worker-class uvicorn.workers.UvicornWorker --bind=0.0.0.0:8100