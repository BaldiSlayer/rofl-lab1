import requests
import json

URL = 'http://localhost:11434/api/generate'
MODEL = 'llama3'


def generate_response(prompt: str) -> str:
    data = {
        "model": MODEL,
        "prompt": prompt,
    }

    response = requests.post(URL, json=data, stream=True)

    response.raise_for_status()

    data_string = ""
    for chunk in response.iter_content(chunk_size=None):
        ans_dict = json.loads(chunk.decode('utf-8'))
        data_string += ans_dict["response"]

    return data_string


print(generate_response("how do you do?"))
