# ---
# jupyter:
#   jupytext:
#     text_representation:
#       extension: .py
#       format_name: percent
#       format_version: '1.3'
#       jupytext_version: 1.14.7
#   kernelspec:
#     display_name: Python 3
#     name: python3
# ---

# %% [markdown] id="6p9cJFLdEfS1"
# ### Гипотезы
# 1. Убрать ненужные слова по типу: расскажи, докажи, из, ...
#
#     Есть два пути: выбрать опорные слова из определений и оставить только их или собрать служебные и второстепенные слова и выкинуть служебные
# 2. На каком языке лучше английском или русском? (кажется, что модель вообще не воспринимает русский и проще дообучить). Что лучше: улучшить перевод или дообучить модель?
# 3. Дообучить модель на корпусте текста (учебник мб)
# 4. Взять другую, русскую модель и использовать

# %% [markdown] id="CDA-IlpYI4vM"
# ### Как считаю метрику?
# Беру два вопроса: пользовательский и из базы знаний, привожу их к одному виду (прописные буквы, начальная форма, удаление стоп слов), потом считаю пересечение и делю на количество слов в вопросе из БЗ.
#
# **Улчушение**: для каждого слова определить его важность: как количество вхождений в БЗ, либо ручками составить словарь с важными терминами.

# %% colab={"base_uri": "https://localhost:8080/"} id="glT_rP-9pK0q" outputId="00326484-75c8-4536-aa96-73a3a2f85f14"
# !pip install pymorphy2

# %% id="PMtjmNHRSAE3"
import re
import nltk
from nltk.tokenize import sent_tokenize, RegexpTokenizer
from nltk.corpus import stopwords
import pymorphy2
from collections import Counter

# %% colab={"base_uri": "https://localhost:8080/"} id="Fa2StkXMo6Ko" outputId="33603799-ae45-4f98-c2ad-b6afe1cfe529"
nltk.download('stopwords')
nltk.download('punkt')
nltk.download('punkt_tab')

stopwords_ru = stopwords.words('russian')

tokenizer = RegexpTokenizer('\w+')

morph = pymorphy2.MorphAnalyzer()


# %% id="TugFS0eVp417"
def make_beautiful(text):
    text = text.lower() # единый формат
    text = re.sub('\n|\t|\r', ' ', text) # убираем лишнее
    sents = sent_tokenize(text) # разбиваем на предложения
    tokenizer = RegexpTokenizer('\w+')
    sents_tokenize = [tokenizer.tokenize(item) for item in sents] # разбиваем на слова
    sents_tokenize = [[item for item in sent if item not in stopwords_ru] for sent in sents_tokenize] # лемматизируйте все слова из датасета
    sents_tokenize = [[morph.normal_forms(item)[0] for item in sent] for sent in sents_tokenize] # лемматизируйте все слова из датасета
    words = [item for sent in sents_tokenize for item in sent] # получаем обработанные слова

    return words


# %% id="1hYwHpE5q6Dt"
def calculate_k(user, database):
    in_both = list(set(user) & set(database))
    k = len(in_both) / len(database)

    return k


# %% id="bt5_a_G4dt6O"
question = make_beautiful('Как можно вычислить частичные производные Антимирова?')

# %% id="jW752ahHdmU0"
user_question = make_beautiful('че такое эти ваши производные антимирова')

# %% id="fGNiUj7KrcNf"
K = calculate_k(user_question, question)

# %% colab={"base_uri": "https://localhost:8080/"} id="5OcicfL0rjWO" outputId="edd7153d-e840-42c9-cdf4-2ceeed729bf7"
K

# %% id="tNm6wd71ITCU"
word_dict = Counter(question)
word_dict.most_common()[:20]
words = word_dict.most_common()

# %% colab={"base_uri": "https://localhost:8080/"} id="NuI4Reh2WkCo" outputId="732bf5b4-bcde-4a1e-ab64-509a4f96f451"
words

# %% id="CYBllffWWuAz"
