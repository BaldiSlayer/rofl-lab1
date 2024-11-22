import requests

URL = 'http://0.0.0.0:8100/'

def get_similar(question: str):
    res = requests.post(URL + '/search_similar', data={"question": question})

    return res.text


if __name__ == "__main__":
    print(get_similar("Регулярен ли язык Дика"))
