from mistralai import Mistral

import app.config.config as config


def create_mistral_client():
    conf = config.SingletonConfig.get_instance()

    return Mistral(api_key=conf.get_mistral_api_key())


mistral_client = create_mistral_client()


def get_chat_response(prompt, context=None, model="open-mistral-7b"):
    """
    Получает ответ от LLM на основе предоставленного промпта и контекста.

    :param model:
    :param prompt: Строка с текстовым промптом от пользователя.
    :param context: Строка с дополнительным контекстом (опционально).
    :return: Ответ от LLM.
    """

    if not mistral_client:
        return "Ошибка: Клиент не инициализирован."

    messages = []
    if context:
        messages.append({"role": "system", "content": context})

    messages.append({"role": "user", "content": prompt})

    chat_response = mistral_client.chat.complete(
        model=model,
        messages=messages,
        n=1,
    )

    return chat_response.choices[0].message.content
