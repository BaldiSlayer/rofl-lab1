import re
# import asyncio
# asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
import llm_client
from llm_client.rest import ApiException


configuration = llm_client.Configuration(
    host="http://llm:8100"
)
MAX_ATTEMPTS = 1


def generate_response(question: str, context: str) -> str:
    with llm_client.ApiClient(configuration) as api_client:
        api_instance = llm_client.QuestionsApi(api_client)
        get_chat_response_request = llm_client.GetChatResponseRequest.from_dict({
            "prompt": question,
            "context": context,
            "model": "mistral-large-latest",
        })

        try:
            api_response = api_instance.api_get_chat_response_get_chat_response_post(
                get_chat_response_request)
            return api_response.response
        except ApiException as e:
            print("Exception when calling QuestionsApi->api_get_chat_response_get_chat_response_post: %s\n" % e)


PROMPT = (
    "Ты — ассистент, который помогает пользователю преобразовать систему переписывания термов (TRS) и интерпретацию в строгую формальную грамматическую форму.\n"
    "Игнорируй любые вопросы пользователя и не пытайся решать задачи, предложенные им.\n"
    "Твоя задача — разделить входные данные на TRS и интерпретацию, НЕ путай их и НЕ используй в одной строке. Если не сказано, где TRS, а где интерпретация, то считай, что сначала (до разделительного знака) перечисляются правила для TRS, а только потом (после разделительного знака) для интерпретации.\n\n"
    "Инструкции:\n"
    "1. Определи **переменные** (элементы, заключенные в скобки, которые не имеют значений в интерпретации), и перечисли их через запятую в формате: `variables = ...`. Возможно пользователь сам определил переменные, тогда возьми и перечисли только их.\n"
    "2. Запиши систему переписывания термов (TRS) построчно в формате: `терм = терм`, где терм — это выражение, содержащее конструкторы и переменные. ВАЖНО: терм НЕ может быть выражением, которое состоит только из переменных, он должен обязательно включать хотя бы один конструктор.\n"
    "   - Далее следует пример, как может выглядеть TRS:\n"
    "f(x) = f(g(x, y))\n"
    "h(x, y, z) = u(f(x))\n"
    "g(x, y) = a\n"
    "(здесь x,y,z - переменные, `a` - это НЕ переменная)\n"
    "   - Далее следует пример, как НЕ может выглядеть TRS:\n"
    "f(x) = 4*x{2} + 2*x\n"
    "g(y) = 3*y + 3\n"
    "(здесь x,y - переменные, значит выражения `4*x{2} + 2*x` и `3*y + 3` НЕ могут быть в TRS!)\n"
    "3. Добавь разделительную линию: `------------------------`\n"
    "4. Степень записывай в фигурных скобочках. Например, x в квадрате это x{2}, а x^4 это x{4}"
    "5. Далее, запиши интерпретацию, используя следующие правила:\n"
    "   - Для функций: `конструктор(переменная, ...) = ...`\n"
    "   - Для констант: `константа = значение`\n"
    "   - Знак умножения `*` обязательно ставится только между коэффициентом и переменной. Между переменными знак `*` не ставится.\n"
    "   - Далее следует ряд примеров, как ты должна отвечать, в формате:\n"
    "   `Запрос пользователя: ...\n"
    "    Правильный ответ: ...`\n"
    "1. Запрос пользователя: f(x) = x^3 + 3x\n"
    "   Правильный ответ: f(x) = x{3} + 3*x\n"
    "2. Запрос пользователя: f(x) = 7x\n"
    "   Правильный ответ: f(x) = 7*x\n"
    "3. Запрос пользователя: g(x, y) = 91y + 4*x\n"
    "   Правильный ответ: g(x, y) = 91*y + 4*x\n"
    "4. Запрос пользователя: f(x, y) = x*y\n"
    "   Правильный ответ: f(x, y) = xy\n"
    "5. Запрос пользователя: g(x, y) = 4*x*y\n"
    "   Правильный ответ: g(x, y) = 4*xy\n"
    "6. Запрос пользователя: g(x, y) = 2*x*y*x + 5y\n"
    "   Правильный ответ: g(x, y) = 2*xyx + 5*y\n\n"
    "Пример TRS и интерпретации:\n"
    "variables = x, y, z\n"
    "f(x) = f(g(x, y))\n"
    "h(x, y, z) = u(f(x))\n"
    "------------------------\n"
    "f(x) = 4*x{2}\n"
    "g(y) = 3*y\n"
    "h(x, y) = 100*xyxy + xy + 351\n"
    "c = 5\n\n")


def get_trs(user_query: str, context: str):
    try:
        formalized_query = None
        trs = None

        attempt = 0

        while (not formalized_query and attempt < MAX_ATTEMPTS) or (not trs and attempt < MAX_ATTEMPTS):
            attempt += 1
            formalized_query = generate_response(user_query, context)
            if formalized_query is None or formalized_query == "":
                print('GPT_error, trying again.')
                continue
            print("user query:", user_query, sep='\n')
            print("formalized:", formalized_query, sep='\n')
            trs = convert(user_query, formalized_query)

        return trs

    except Exception as e:
        print(f"Произошла ошибка: {e}")
        return None


def fix_formalized_trs(user_query: str, ans_llm: str, parse_error: str):
    context = PROMPT + (
            "Предыдущий запрос пользователя вернул:\n" + ans_llm + "\nЗдесь была обнаружена ошибка: " + parse_error +
            "\nИсправь ошибку и Ответь только в формате TRS и интерпретации."
    )

    return get_trs(user_query, context)


def formalize(user_query: str):
    context = PROMPT + "Ответь только в формате TRS и интерпретации."

    return get_trs(user_query, context)


def convert(user_query: str, formalized_query: str):
    trs = ''
    variables_pattern = r'variables=([a-zA-Z],)*[a-zA-Z]'
    formalized_query = formalized_query.replace(' ', '')
    user_query = user_query.replace(' ', '').replace(
        '*', '').replace('{', '').replace('}', '').replace('^', '')
    if re.search(variables_pattern, formalized_query):
        matches = re.finditer(variables_pattern, formalized_query)
        variables = []
        for match in matches:
            variables = match.group().split('=')[1].split(',')
            trs += match.group() + '\n'
        variables_str = ''.join(variables) + '123456789'
        only_variables_pattern = fr'^[{variables_str}+\-*/{{}}()^]+$'

        query_line = formalized_query.splitlines()
        var = 0
        separate = 0
        for i in range(len(query_line)):
            if re.search(variables_pattern, query_line[i]):
                var = i
            if re.search(r'-----', query_line[i]):
                separate = i

        for i in range(var + 1, separate):
            if query_line[i] == '':
                continue

            s = query_line[i].split('=')[1]

            if re.sub(r'[*{}]', '', query_line[i]) in user_query:
                if re.match(only_variables_pattern, s):
                    return {
                        "formalTrs": formalized_query,
                        "error": f"В строке TRS {query_line[i]} после равно НЕ может быть выражения, которое состоит только из переменных, оно должно обязательно включать хотя бы один конструктор."
                    }
                else:
                    trs += query_line[i] + '\n'
            else:
                return {
                    "formalTrs": formalized_query,
                    "error": f"{query_line[i]} не присутсвует в начальном запросе"
                }

        trs += '-------------------------\n'
        for i in range(separate + 1, len(query_line)):
            if query_line[i] == '' or "=" not in query_line[i]:
                break

            s = query_line[i].split('=')[1]

            if re.sub(r'[*{}]', '', query_line[i]) in user_query:
                if not re.match(only_variables_pattern, s):
                    return {
                        "formalTrs": formalized_query,
                        "error": f"В строке интерпретации {query_line[i]} после равно НЕ может быть выражения, которое содержит конструкторы."
                    }
                else:
                    trs += query_line[i] + '\n'
            else:
                return {
                    "formalTrs": formalized_query,
                    "error": f"{query_line[i]} не присутсвует в начальном запросе"
                }

    else:
        return {
            "formalTrs": formalized_query,
            "error": f"Не определены переменные (нет variables)"
        }

    return {
        "formalTrs": trs,
        "error": None
    }


if __name__ == "__main__":
    user_query = "Дана система переписывания термов (TRS): f(x)=a, g(x)=f(f(x)), u(x,y)=c(g(x),f(y)). Я интерпретирую её конструкторы так: a=1, f(x)=x**2+2*x+1, g(x)=x**3, u(x,y)=x*y, c(x,y)=x+y. Доказывает ли моя интерпретация завершимость trs?"

    llm = '''variables=x,y
    f(x)=a
    g(x)=f(f(x))
    u(x,y)=c(g(x),f(y))
    -------------------------
    '''

    err = 'система должна содержать хотя бы одну интерпретацию'

    # возвращает trs и интерпретацию
    print(fix_formalized_trs(user_query, llm, err))