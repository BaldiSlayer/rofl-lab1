FROM public.ecr.aws/docker/library/python:3.9-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH="/LLM"

RUN apt update && apt install -y curl

COPY /LLM/requirements.txt /LLM/requirements.txt
RUN pip install --no-cache-dir -r /LLM/requirements.txt

COPY /LLM /LLM
COPY /data/data.yaml /LLM/app/init_embeddings/data.yaml

WORKDIR /LLM

RUN python /LLM/app/init_embeddings/init_embeddings.py /LLM/app/init_embeddings/data.yaml /LLM/vectorized_data.faiss && \
    mv /LLM/app/init_embeddings/data.yaml /LLM/data.yaml

COPY /similarity-tester /similarity-tester

RUN chmod +x /similarity-tester/start.sh

ENTRYPOINT ["/bin/bash", "/similarity-tester/start.sh"]
