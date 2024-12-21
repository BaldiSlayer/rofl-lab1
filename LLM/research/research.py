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
# 2. На каком языке лучше (хотим русский)
# 3. Дообучить русскую нейронку (недавно вышла Т-pro)

# %% [markdown] id="eBdqyMMU0_GZ"
# ## качаю нужные библиотеки

# %% colab={"base_uri": "https://localhost:8080/"} id="8nxzUr8-3NsH" outputId="0438dc07-03b1-46d9-be9d-e352563190d9"
# !pip install pymorphy2

# %% colab={"base_uri": "https://localhost:8080/"} id="E50FJEIezntI" outputId="08a7a992-129e-417e-db52-74076effc814"
# !pip install pdfminer.six

# %% colab={"base_uri": "https://localhost:8080/"} id="SSsLEe_zznq9" outputId="a56346db-5d2b-497b-8e54-eab5466a41d9"
# !pip install 'pdfminer.six[image]'

# %% [markdown] id="CyaR7pFa1Jij"
# ## всякие импорты

# %% id="PMtjmNHRSAE3"
import re # регулярки
import nltk # для стоп-слов
from nltk.tokenize import sent_tokenize # для токенизации
from nltk.tokenize import RegexpTokenizer # разбиваем на слова
from nltk.corpus import stopwords # убираем стоп слова
import pymorphy2 # для лемантизации
from collections import Counter # подсчёт частых
from gensim.models.phrases import Phrases, Phraser # дли биграмм

from gensim.models.word2vec import Word2Vec

from pdfminer.high_level import extract_text # pdf в текст

from google.colab import drive # подключение к диску

# %% colab={"base_uri": "https://localhost:8080/"} id="pe5CpIvY0RY7" outputId="becca389-1a1a-4fee-f6be-7f6c51260b8e"
drive.mount('/content/drive', force_remount=True)

# %% colab={"base_uri": "https://localhost:8080/"} id="1jjMPO7aSlRS" outputId="8efedf04-17d8-49fd-8951-1cea2baa1e4c"
nltk.download('stopwords')
nltk.download('punkt')
nltk.download('punkt_tab')

# %% id="dQ6xQf3A3dSf"
stopwords_ru = stopwords.words('russian')
tokenizer = RegexpTokenizer('\w+')
morph = pymorphy2.MorphAnalyzer()

# %% [markdown] id="p-KjpkXL2vZh"
# ## начинаем работу с текстами

# %% [markdown] id="Ot4HjDfUDErj"
# ### книга

# %% id="VvKlY3SokHGK"
# лекции за 2023 год + материалы итмо + книга хопкрофта на русском
all_texts = []

for i in range(1, 68):
    t = extract_text(f"/content/drive/MyDrive/texts/{i}.pdf")
    # глазами зацепилась за не оч хорошие переносы, поэтому фикс
    t = t.replace('-\n', '')
    all_texts.append(t)


# %% [markdown] id="2iuzfpEoDHi9"
# ## код

# %% [markdown] id="CxCQYjIPsd65"
# ### какие-то мои функции

# %% id="rg8_EDpD8yhb"
# подсчёт метрики k
def my_metrics(user, db, clean_dict=None, trash_dict=None):
    user = make_clean(user)
    db = make_clean(db)

    user = get_words_only(user)
    db = get_words_only(db)

    if clean_dict is not None:
        user = [w for w in user if w in my_dict]
        db = [w for w in db if w in my_dict]

    if trash_dict is not None:
        user = [w for w in user if w not in trash_dict]
        db = [w for w in db if w not in trash_dict]

    print(user)
    print(db)

    in_both = list(set(user) & set(db))

    return len(in_both) / (len(db) + epsilon)


# %% id="LhYu7mWu9wSF"
# берём только слова из предложений
def get_words_only(sentences_tokenize):
    words = [item for sent in sentences_tokenize for item in sent]

    return words


# %% id="RFvD0NeL2Wvu"
def make_clean(text):
    # чистим текст
    text = text.lower()
    text = re.sub('\n|\t|\r', ' ', text)
    # преобразуем в предложения
    sentences = sent_tokenize(text)
    # бъём предложения на слова
    sentences_tokenize = [tokenizer.tokenize(item) for item in sentences]
    # убираем стоп слова
    sentences_tokenize = [[item for item in sent if (item not in stopwords_ru and re.match('^[а-я]*$', item))] for sent in sentences_tokenize]
    # лемантизируем слова
    sentences_tokenize = [[morph.normal_forms(item)[0] for item in sent] for sent in sentences_tokenize]

    return sentences_tokenize


# %% id="QCoWUMKk30Kf"
# словарь популярных слов
def get_dict_words(sentences_tokenize):
    words = [item for sent in  sentences_tokenize for item in sent]
    word_dict = Counter(words)
    common = word_dict.most_common()

    return common


# %% id="rtzmzqEu8iQh"
# полезные слова = часто встречаем их
def get_useful(common):
    my_dict = {}

    for i, j in common:
        if i != 'ε' and len(i) > 2 and re.match('^[а-я]*$', i):
            my_dict[i] = j

    return my_dict


# %% [markdown] id="xLyXh4KBsh7A"
# ### смотрим на текст

# %% id="tnfIQEDcpTXn"
epsilon = 0.001

all_sentences = []
all_words = []

my_dict = {}

# мусорные конструкции
trash_dict = ['че', 'доказать', 'рассказать', 'ваш', 'привести', 'сделать', 'описать', 'дать', 'что']

# %% id="ntRI3jSxo-Xi"
for text in all_texts:
    t = make_clean(text)
    all_sentences = all_sentences + t

    # частые слова
    common = get_dict_words(t)
    all_words += [i for (i, j) in common]
    # «полезные» слова
    useful_dict = get_useful(common)

    new_dict = dict(list(my_dict.items()) + list(useful_dict.items()))
    my_dict = new_dict.copy()

# %% [markdown] id="DqvEnLZ9spVQ"
# ### предобработка перед w2v

# %% id="jSfuYLAAstDw"
bigram = Phrases(all_sentences)
bigram_transformer = Phraser(bigram)

# генератор текстов с биграммами
def bigram_generator():
    for text in all_sentences:
        yield bigram_transformer[[word for word in text]]


# %% id="SLWhw_oRvJUp"
trigram = Phrases(bigram_generator())
trigram_transformer = Phraser(trigram)

def trigram_generator():
    for text in all_sentences:
        yield trigram_transformer[bigram_transformer[[word for word in text]]]


# %% id="_bXigcMdvQLb"
words = [i for i in trigram_generator()]

# %% id="fU_fJ8wnvXDt"
model = Word2Vec(vector_size=300, window=7, min_count=3)
model.build_vocab(words)

# %% colab={"base_uri": "https://localhost:8080/"} id="SOR_qrPIvdNd" outputId="4fd22705-115d-4231-971d-2b7460008482"
model.train(words, total_examples=model.corpus_count, epochs=5)

# %% colab={"base_uri": "https://localhost:8080/"} id="N1Up3Fc31jjO" outputId="6cc8392b-cf95-4889-b80a-cf68b1b0819b"
model.wv.similarity('слово', 'цепочка')

# %% id="T2QHea_0vgBl"
# сохраню потом
model.save('wv_model')

# %% id="irA_1fUlviOa"
trigram_transformer.save('wv_trigramm')

# %% [markdown] id="Hc3EmqlPrQjB"
# # подвал

# %% id="yvs_UWnh2WsV"
# чистим текст из пдфки
test = make_clean(all_texts[0])

# %% id="mZhO3Buy4Sis"
# словарь с частыми словами
common = get_dict_words(test)

# %% id="ubMP3X-2-ixl"
# словарь с полезными терминами
my_dict = get_useful(common)

# %% id="jW752ahHdmU0"
s1 = 'че такое эти ваши производные антимирова'

s2 = 'Как можно вычислить частичные производные Антимирова?'

# %% colab={"base_uri": "https://localhost:8080/"} id="pgytyVDsd6LS" outputId="7f320050-d35a-4d57-e808-2e3489606526"
my_metrics(s1, s2, trash_dict=trash_dict)

# %% [markdown] id="q6QZNMrgDqvw"
# и тут для меня дошло, что, если расписывать вопросы в БЗ, то деление на длину слов вопроса из БЗ (коих будет много) не будет явно и точно отражать близость, но относительно такой штуки можно смотреть динамику
