FROM public.ecr.aws/docker/library/python:3.12-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN apt update && apt install -y curl

COPY /formalize_convert/requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY /formalize_convert .

CMD gunicorn 'openapi_server.__main__:app' --workers 8 --bind=0.0.0.0:8081 --timeout 240
