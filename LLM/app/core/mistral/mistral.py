from mistralai import Mistral

import app.config.config as config


class SingletonTextTranslator:
    _instance = None

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            conf = config.SingletonConfig.get_instance()

            cls._instance = Mistral(api_key=conf.get_mistral_api_key())

        return cls._instance


def get_chat_response(prompt, context=None, model="open-mistral-7b"):
    """
    Получает ответ от LLM на основе предоставленного промпта и контекста.

    :param model:
    :param prompt: Строка с текстовым промптом от пользователя.
    :param context: Строка с дополнительным контекстом (опционально).
    :return: Ответ от LLM.
    """
    client = SingletonTextTranslator.get_instance()

    if not client:
        return "Ошибка: Клиент не инициализирован."

    messages = []
    if context:
        messages.append({"role": "system", "content": context})

    messages.append({"role": "user", "content": prompt})

    chat_response = client.chat.complete(
        model=model,
        messages=messages
        n=1
    )

    return chat_response.choices[0].message.content
