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

# %% colab={"base_uri": "https://localhost:8080/"} id="8nxzUr8-3NsH" outputId="3695f324-f2b5-4762-de99-a9abf732445e"
# !pip install pymorphy2

# %% colab={"base_uri": "https://localhost:8080/"} id="E50FJEIezntI" outputId="303d3cf4-aa08-43da-8ae6-3c5cbd4ad45d"
# !pip install pdfminer.six

# %% colab={"base_uri": "https://localhost:8080/"} id="SSsLEe_zznq9" outputId="da9f2505-0d7d-4a1b-8f59-2f519b076990"
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
from scipy import spatial

from gensim.models.word2vec import Word2Vec
from gensim.models.doc2vec import Doc2Vec, TaggedDocument

from pdfminer.high_level import extract_text # pdf в текст

from google.colab import drive # подключение к диску

# %% colab={"base_uri": "https://localhost:8080/"} id="pe5CpIvY0RY7" outputId="ba38775a-d3bb-4fc6-dd12-2ab9f2eacace"
drive.mount('/content/drive', force_remount=True)

# %% colab={"base_uri": "https://localhost:8080/"} id="1jjMPO7aSlRS" outputId="92d8c29f-0c5c-4069-fcb0-e15ccc8a7752"
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

# %% colab={"background_save": true} id="VvKlY3SokHGK"
# лекции за 2023 год + материалы итмо + книга хопкрофта на русском
all_texts = []

for i in range(1, 69):
    t = extract_text(f"/content/drive/MyDrive/texts/{i}.pdf")
    # глазами зацепилась за не оч хорошие переносы, поэтому фикс
    t = t.replace('-\n', '').replace('\n\n', '\n')
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
    if len(t) > 0:
        all_sentences = all_sentences + t

    # частые слова
    common = get_dict_words(t)
    all_words += [i for (i, j) in common]
    # «полезные» слова
    useful_dict = get_useful(common)

    new_dict = dict(list(my_dict.items()) + list(useful_dict.items()))
    my_dict = new_dict.copy()

# %% [markdown] id="3H7KF6QYK_TE"
# ### предобработка перед doc2vec

# %% id="EpWP_U96LFH4"
vocab_index = [TaggedDocument(text,[k]) for k, text in enumerate(all_sentences)]

model = Doc2Vec(vector_size=300, window=1, min_count=3, workers=6, epochs=20)

model.build_vocab(vocab_index)
model.train(vocab_index, total_examples=model.corpus_count, epochs=20)

# %% colab={"base_uri": "https://localhost:8080/"} id="vtVaalz1Mg4i" outputId="517ce60a-d566-486a-fbfb-cb85af5ff0b1"
s1 = model.infer_vector('дка'.split())
s2 = model.infer_vector('регулярные языки и их свойства'.split())

cos_distance = spatial.distance.cosine(s1, s2)

cos_distance

# %% colab={"base_uri": "https://localhost:8080/"} id="utLC_GHmN8SS" outputId="821f3ee5-9c89-4957-b23a-df67bc53f534"
similar_vec = model.docvecs.most_similar([s2], topn=10)

similar_sent = [" ".join(vocab_index[top[0]].words) for top in similar_vec]

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


# %% id="QR-OKoYAPPqO"
ok = [i for i in trigram_generator()]

# %% id="fU_fJ8wnvXDt"
model = Word2Vec(vector_size=300, window=10, min_count=3, workers=4)
model.build_vocab(ok)

# %% colab={"base_uri": "https://localhost:8080/"} id="moDBxGzkPx8U" outputId="2757a0cc-6c05-497b-ac3b-80957175ff4f"
model.corpus_count

# %% colab={"base_uri": "https://localhost:8080/"} id="SOR_qrPIvdNd" outputId="a4f6658a-222a-4896-edb8-75e071fc8276"
model.train(ok, total_examples=model.corpus_count, epochs=500)

# %% colab={"base_uri": "https://localhost:8080/"} id="N1Up3Fc31jjO" outputId="b0e92a6d-da09-4b0c-b743-9184e332dfa8"
model.wv.most_similar('дка')

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
