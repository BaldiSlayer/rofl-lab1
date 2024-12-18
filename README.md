# rofl-lab1

## Гайд на отправку запросов к LLM

### Пререквизиты
- склонирован репозиторий
- установлен [Docker](https://www.docker.com)

1) Собираем `Docker` образ:
    ```bash
    docker build -f ./dockerfiles/llm.dockerfile -t my-t .
    ```
2) Запускаем `Docker` контейнер
    ```bash
   # подставить API ключ mistral
    docker run -p 8100:8100 -e MISTRAL_API_KEY=<API_КЛЮЧ_MISTRAL> my-t
    ```
3) Переходим на http://localhost:8100/docs. Это раздел с документацией к `API`, отсюда можно отправлять
запросы на это самое `API`. Делается это так: выбираем любую ручку (например `/ping`), кликаем по ней, там
будет кнопка `Try it out`, нажимаем туда - появляется поле ввода для всех параметров, вводим и кликаем на кнопку
`Execute`. У каждой ручки также есть описание, что она делает. Для поиска похожих для вопроса надо пользоваться
[/search_similar](http://0.0.0.0:8100/docs#/Questions/api_search_similar_search_similar_post), а для отправки запроса
в `mistral` - [/get_chat_response](http://0.0.0.0:8100/docs#/Questions/api_get_chat_response_get_chat_response_post).

**Важно**: сейчас в проде используем `mistral-large-2411`. Эта же версия
используется по умолчанию в ручке `/get_chat_response`.

## Инструкция по настройке Telegram-бота

### 1. Создание бота в Telegram
1. Запустите бота [@BotFather](https://t.me/BotFather).
2. Начните с ним чат и используйте команду `/newbot`.
3. Следуйте инструкциям для создания бота и получите токен.

### 2. Клонирование репозитория
Склонируйте репозиторий с исходным кодом.

### 3. Создание файла .env
В корне репозитория создайте файл `.env` и заполните его следующим образом:

```
MISTRAL_API_KEY=<mistal token>
TGTOKEN=<token>
POSTGRES_PASSWORD=strong
POSTGRES_USER=bb
POSTGRES_DB=some
GHTOKEN=<ask @Baldislayer for token>
```

### 4. Запуск через Docker Compose
Для сборки и запуска контейнеров выполните:

```bash
docker compose up --build
```

### 5. Подключение к PostgreSQL
После запуска контейнера подключитесь к контейнeру PostgreSQL с помощью команды:

```bash
docker exec -it postgres /bin/sh
```

И выполните команду:
```bash
psql -h localhost -d some -U bb
```

### 6. Применение миграций
Выполните для применения миграций:

```sql
CREATE SCHEMA tfllab1;

CREATE TABLE tfllab1.user_state (
       user_id BIGINT PRIMARY KEY,
       state INT NOT NULL,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE tfllab1.extraction_result (
       user_id BIGINT PRIMARY KEY,
       user_request TEXT,
       formalize_result TEXT,
       parse_result JSONB,
       parse_error TEXT
);

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_user_state_updated_at BEFORE UPDATE ON tfllab1.user_state
       FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();

CREATE INDEX CONCURRENTLY user_state_updated_at_index ON tfllab1.user_state (updated_at);

CREATE TABLE tfllab1.user_lock (
       user_id BIGINT PRIMARY KEY,
       expires_at TIMESTAMP NOT NULL,
       instance_id TEXT NOT NULL
);

CREATE INDEX user_lock_expires_at_index ON tfllab1.user_lock (expires_at);
CREATE INDEX user_lock_instance_id_index ON tfllab1.user_lock (instance_id);
```

8. Готово! 🎉
