FROM python:3.12-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

COPY /interpret/requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY /interpret .

CMD gunicorn 'openapi_server.__main__:app' --workers 8 --bind=0.0.0.0:8081
