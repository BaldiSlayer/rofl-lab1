import requests
import json

URL = 'http://localhost:11434/api/generate'
MODEL = 'llama3.1'


def generate_response(prompt: str) -> str:
    data = {
        "model": MODEL,
        "prompt": prompt,
    }

    response = requests.post(URL, json=data, stream=True)

    response.raise_for_status()

    data_string = ""
    for chunk in response.iter_content(chunk_size=None):
        print(chunk)
        ans_dict = json.loads(chunk.decode('utf-8'))
        data_string += ans_dict["response"]

    return data_string


print(generate_response("Как дела?"))


# git pull https://github.com/BaldiSlayer/rofl-lab1.git formalize_convert