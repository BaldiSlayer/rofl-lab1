#!/bin/sh

echo "initializing embeddings..."

export PYTHONPATH="/LLM"

python3 /LLM/app/utils/create_embeddings.py /LLM/app/utils/data.yaml /LLM/vectorized_data.faiss

mv /LLM/app/utils/data.yaml /LLM/data.yaml

echo "embeddings were successfully initialized..."

exec gunicorn app.main:app --workers 1 --worker-class uvicorn.workers.UvicornWorker --bind=0.0.0.0:8100 --timeout 240
