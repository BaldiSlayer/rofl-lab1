from mistralai import Mistral
import os

from ...utils.Mistral.config import api_key


# from getAPIkey import make_api_key

def make_api_key():
    """
    place your API key in environment
    """
    os.environ["MISTRAL_API_KEY"] = api_key


def initialize_client():
    """Инициализирует клиента Mistral с API-ключом."""
    make_api_key()
    api_key = os.environ.get("MISTRAL_API_KEY")
    if api_key:
        print("API key найден")
        return Mistral(api_key=api_key)
    else:
        print("API key не найден.")
        return None


def get_chat_response(prompt, context=None, model="open-mistral-7b"):
    """
    Получает ответ от LLM на основе предоставленного промпта и контекста.

    :param prompt: Строка с текстовым промптом от пользователя.
    :param context: Строка с дополнительным контекстом (опционально).
    :return: Ответ от LLM.
    """
    client = initialize_client()
    if not client:
        return "Ошибка: Клиент не инициализирован."

    messages = []
    if context:
        messages.append({"role": "system", "content": context})

    messages.append({"role": "user", "content": prompt})

    chat_response = client.chat.complete(
        model=model,
        messages=messages
    )
    return chat_response.choices[0].message.content


# prompt = "напиши предыдущее сообщение"
# context = "это предыдущее сообщение: ахахахахах"
# response = get_chat_response(prompt, context)
# print(response)
