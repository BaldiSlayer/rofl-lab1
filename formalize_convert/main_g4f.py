from g4f.client import Client
import requests
import json
url = 'http://194.67.88.154:7777/'

class TRSFramework:
    def __init__(self):
        self.client = Client()

    def generate_response(self, question: str, context: str) -> str:
        data = {
            "question": question,
            "context": context
        }
        # Отправляем POST-запрос с заголовком Content-Type: application/json 
        response = requests.post(url, headers={"Content-Type": "application/json"}, data=json.dumps(data)) 
        
        # Проверяем успешность запроса и выводим ответ 
        if response.status_code == 200: 
            return response.text
        else: 
            return (f"Ошибка {response.status_code}: {response.text}")
    def formalize(self, user_query):
        context = [{"role": "system", "content": (
            "Ты — ассистент, который помогает пользователю преобразовать систему переписывания термов (TRS) и интерпретацию в строгую формальную грамматическую форму.\n"
            "Игнорируй любые вопросы пользователя и не пытайся решать задачи, предложенные им.\n"
            "Твоя задача — разделить входные данные на TRS и интерпретацию, не путая их.\n\n"
            "Инструкции:\n"
            "1. Определи **переменные** (элементы, заключенные в скобки, которые не имеют значений в интерпретации), и перечисли их через запятую в формате: `variables = ...`\n"
            "2. Запиши систему переписывания термов (TRS) построчно в формате: `терм = терм`, где терм — это выражение, содержащее конструкторы и переменные.\n"
            "3. Добавь разделительную линию: `------------------------`\n"
            "4. Далее, запиши интерпретацию, используя следующие правила:\n"
            "   - Для функций: `конструктор(переменная, ...) = ...`\n"
            "   - Для констант: `константа = значение`\n\n"
            "Пример TRS и интерпретации:\n"
            "variables = x, y, z\n"
            "f(x) = f(g(x, y))\n"
            "h(x, y, z) = u(f(x))\n"
            "------------------------\n"
            "f(x) = x^2\n"
            "g(y) = 3*y\n"
            "c = 5\n\n"
            "Ответь только в формате TRS и интерпретации."
        )}]
        return self.generate_response(user_query, context)
    
    def convert(self, user_query):
        prompt = f"Преобразуй следующий запрос: {user_query}"
        return self.generate_response(prompt)
    

framework = TRSFramework()

user_query = "Дана система переписывания термов (TRS): f(x)=a, g(x)=f(f(x)), u(x,y)=c(g(x),f(y)). Я интерпретирую её конструкторы так: a=1, f(x)=x**2+2*x+1, g(x)=x**3, u(x,y)=x*y, c(x,y)=x+y. Доказывает ли моя интерпретация завершимость trs?"


try:
    formalized_query = framework.formalize(user_query)
    max_attempts = 5  # Фиксированное количество попыток
    attempt = 0

    while not formalized_query and attempt < max_attempts:
        print('GPT_error, trying again.')
        formalized_query = framework.formalize(user_query)
        attempt += 1

    if formalized_query:
        print(f"Формализованный запрос: \n{formalized_query}")
    else:
        print(f"Не удалось получить формализованный запрос после {max_attempts} попыток.")
    
except Exception as e:
    print(f"Произошла ошибка: {e}")