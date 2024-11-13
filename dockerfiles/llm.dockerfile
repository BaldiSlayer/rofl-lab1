FROM public.ecr.aws/docker/library/python:3.9-slim AS create-embeddings

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

RUN apt update && apt install -y curl

COPY /LLM/requirements.txt /LLM/requirements.txt
RUN pip install --no-cache-dir -r /LLM/requirements.txt

COPY /LLM /LLM
COPY /data/data.yaml /LLM/app/utils/data.yaml

WORKDIR /LLM

RUN chmod +x /LLM/app/utils/entrypoint.sh

ENTRYPOINT ["/LLM/app/utils/entrypoint.sh"]