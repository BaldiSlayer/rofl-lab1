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
POSTGRES_USER=admin
POSTGRES_DB=dbtfl
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
psql -h localhost -d dbtfl -U admin
```

### 6. Применение миграций
Выполните [следующие](https://github.com/BaldiSlayer/rofl-lab1/blob/main/postgresql/final_schema.sql) для применения миграций (первую строчку копировать не надо).

### 7. Готово! 🎉
