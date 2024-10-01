import requests
import json

class TRSFramework:
    def __init__(self, server_url):
        self.server_url = server_url  # URL для вашего локального сервера LLM

    def generate_response(self, prompt: str) -> str:
        data = {
            "model": "llama3.1",
            "prompt": prompt,
        }

        response = requests.post(self.server_url, json=data, stream=True)

        response.raise_for_status()

        data_string = ""
        for chunk in response.iter_content(chunk_size=None):
            ans_dict = json.loads(chunk.decode('utf-8'))
            data_string += ans_dict.get("response", "")

        return data_string

    def formalize(self, user_query):
        prompt = (
            f"Исправь опечатки и грамматические ошибки в следующем тексте: '{user_query}'. "
            "Выдай только исправленный текст, без дополнительных комментариев."
        )
        
        return self.generate_response(prompt)
    
    def convert(self, user_query):
        prompt = (
            "prompt"
        )

        return self.generate_response(prompt)
    
    def clear_mem(self):
        prompt = ("/clear")
        return self.generate_response(prompt)

server_url = 'http://localhost:11434/api/generate'

framework = TRSFramework(server_url)

user_query = "Дана сстема прписыванмя трмов: S-> aS, S -> b. Я интерпретирую её конструкторы как криволинейные функции. Доказывает ли моя интерпретация завершимость trs?"

try:
    # Очистка памяти чата
    clear_mem_status = framework.clear_mem()

    formalized_query = framework.formalize(user_query)
    print(f"Формализованный запрос: {formalized_query}")
    
except Exception as e:
    print(f"Произошла ошибка: {e}")