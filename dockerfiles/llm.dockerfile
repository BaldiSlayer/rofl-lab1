FROM public.ecr.aws/docker/library/python:3.9-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH="/LLM"

RUN apt update && apt install -y curl

COPY /LLM/requirements.txt /LLM/requirements.txt
RUN pip install --no-cache-dir -r /LLM/requirements.txt

COPY /LLM/app /LLM/app
COPY /LLM/__init__.py /LLM/__init__.py

COPY /data/data.yaml /LLM/app/init_embeddings/data.yaml

WORKDIR /LLM

RUN python /LLM/app/init_embeddings/init_embeddings.py /LLM/app/init_embeddings/data.yaml /LLM/vectorized_data.faiss && \
    mv /LLM/app/init_embeddings/data.yaml /LLM/data.yaml

CMD gunicorn app.main:app --workers 8 --worker-class uvicorn.workers.UvicornWorker --bind=0.0.0.0:8100 --timeout 240
