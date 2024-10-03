import requests, json

URL = 'http://localhost:8100/get_chat_response'
MODEL = 'open-mistral-7b'


def generate_response(prompt: str, context: str) -> str:
    data = {
        "prompt": prompt,
        "context": context,
        "model": MODEL
    }

    response = requests.post(URL, json=data, stream=True)

    response.raise_for_status()

    data_string = ""
    for chunk in response.iter_content(chunk_size=None):
        ans_dict = json.loads(chunk.decode('utf-8'))
        data_string += ans_dict["response"]

    return data_string


print(generate_response("how do you do?", "norm"))

