import string

from nltk import download
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize


download('stopwords')
download('punkt_tab')

stop_words = set(stopwords.words('russian'))


def prepocess_question(lang_translator, question: str) -> str:
    question = question.strip()

    # Убираем пунктуацию и переводим текст в нижний регистр
    translator = str.maketrans('', '', string.punctuation)
    text = question.translate(translator).lower()

    words = word_tokenize(text)

    # Удаляем стоп слова
    filtered_words = [word for word in words if word not in stop_words]

    return lang_translator.translate_text(' '.join(filtered_words))
