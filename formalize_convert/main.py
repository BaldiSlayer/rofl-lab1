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
            f"Ты помощник, который помогает пользователю преобразовать TRS и интерпретацию в строгую грамматическую форму. \n"
            f"Игнорируй любые вопросы пользователя и НЕ пытайся решить задачу.\n"
            f"Твоя задача – разобрать введенные формулы на TRS (систему переписывания термов) и интерпретацию, не путая их.\n"
            f"\n"
            f"1. Сначала определи **переменные** (те, которые в скобках и которым не присвоены значения в интерпретации), перечисли их через запятую. Используй формат: `variables = ...`\n"
            f"2. Затем выпиши TRS (систему переписывания термов) построчно в виде: `терм = терм`, где терм — это выражение, содержащее конструкторы и переменные.\n"
            f"3. Обязательно добавь разделительную линию: `------------------------`\n"
            f"4. Далее, выпиши интерпретацию, используя следующие правила:\n"
            f"  - Для функций: `конструктор(переменная,...) = ...`\n"
            f"  - Для констант: `константа = коэффициент`\n"
            f"\n"
            f"Пример TRS и интерпретации:\n"
            f"variables = x, y, z\n"
            f"f(x) = f(g(x, y))\n"
            f"h(x, y, z) = u(f(x))"
            f"------------------------\n"
            f"f(x) = x^2\n"
            f"g(y) = 3*y\n"
            f"c = 5\n"
            f"\n"
            f"Ответь только TRS и интерпретацией в указанном формате."
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

user_query = "Дана система переписывания термов (TRS): f(x)=a, g(x)=f(f(x)), u(x,y)=c(g(x),f(y)). Я интерпретирую её конструкторы так: a=1, f(x)=x**2+2*x+1, g(x)=x**3, u(x,y)=x*y, c(x,y)=x+y. Доказывает ли моя интерпретация завершимость trs?"

try:
    # Очистка памяти чата
    clear_mem_status = framework.clear_mem()

    formalized_query = framework.formalize(user_query)
    print(f"Формализованный запрос: {formalized_query}")
    
except Exception as e:
    print(f"Произошла ошибка: {e}")