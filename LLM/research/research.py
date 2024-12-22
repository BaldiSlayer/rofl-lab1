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

# %% colab={"base_uri": "https://localhost:8080/"} id="8nxzUr8-3NsH" outputId="03bfd5b3-9566-4742-cd49-d58b218c35e8"
# !pip install pymorphy2

# %% colab={"base_uri": "https://localhost:8080/"} id="E50FJEIezntI" outputId="b879c8b8-4c38-4678-9106-742e3a7ff19c"
# !pip install pdfminer.six

# %% colab={"base_uri": "https://localhost:8080/"} id="SSsLEe_zznq9" outputId="83671d3b-5276-4067-a625-d84fdc89d09c"
# !pip install 'pdfminer.six[image]'

# %% colab={"base_uri": "https://localhost:8080/"} id="2oObxkmlGca2" outputId="1e9711d3-9e44-4a96-f970-10ad653fddf6"
# !pip install faiss

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
import numpy as np
import tensorflow as tf

from gensim.models.word2vec import Word2Vec
from gensim.models.doc2vec import Doc2Vec, TaggedDocument

from pdfminer.high_level import extract_text # pdf в текст

from google.colab import drive # подключение к диску

# %% colab={"base_uri": "https://localhost:8080/"} id="pe5CpIvY0RY7" outputId="2ec5e0cb-2628-4afc-a03f-2022d66c3596"
drive.mount('/content/drive', force_remount=True)

# %% colab={"base_uri": "https://localhost:8080/"} id="1jjMPO7aSlRS" outputId="2db0b24a-f7dd-4752-8f36-8b9d10edd7b9"
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

# %% colab={"base_uri": "https://localhost:8080/"} id="vtVaalz1Mg4i" outputId="65fbfa5c-44d7-4fc6-ae8d-19ec324b18b8"
s1 = model.infer_vector('дка'.split())
s2 = model.infer_vector('регулярные языки и их свойства'.split())

s1 = tf.math.l2_normalize(s1).numpy()
s2 = tf.math.l2_normalize(s2).numpy()

cos_distance = spatial.distance.cosine(s1, s2)

cos_distance

# %% colab={"base_uri": "https://localhost:8080/"} id="utLC_GHmN8SS" outputId="81e0c3b2-39a9-47cb-e72b-c791f16af0ba"
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

# %% colab={"base_uri": "https://localhost:8080/"} id="3C5whEvNwxFP" outputId="e80de487-b616-43bc-ec08-4ac435f4be67"
# !pip install FlagEmbedding

# %% colab={"base_uri": "https://localhost:8080/", "height": 49, "referenced_widgets": ["f82667ff6c1942e4952121ba5e1923ea", "344bdcd2123045eb83f34c54f9858a13", "b5d7daf0a9f54f35a658a41ed9f5cec8", "5af407bf6155464d8de2504fc618b4fd", "84d5c47c9fd94cc1a5c15b17629ceab6", "289845da237546a1a5bb6f6ce6d99154", "a09ee8e503a8484e8b28d13b36b224fa", "0ca3176a2c764ff082044b80916cec27", "c23595416e8d46eaa4ec7195fee48838", "8cde040e5d49452997d4de90649f486a", "078cd8bd952e4a438b458cea57602748"]} id="VI7BNclkw9Ou" outputId="8915df4e-846b-4b8e-9ab2-d212e61e8ae4"
from FlagEmbedding import BGEM3FlagModel

model = BGEM3FlagModel('BAAI/bge-m3',  use_fp16=True)

# %% colab={"base_uri": "https://localhost:8080/"} id="M6pyw0cizrKr" outputId="dd24e2a1-30c2-4aee-a733-2eb17da9dd2f"
from sentence_transformers import SentenceTransformer

mmm = SentenceTransformer("sentence-transformers/all-mpnet-base-v2")

# %% id="vaQbzNkS2SG9"
sentences_1 = ['че такое эти ваши производные антимирова',
               "великая и всемогущая модель, расскажи мне про регулярные языки",
               "правда, что любой язык детерминирован?"]
sentences_2 = ['Как можно вычислить частичные производные Антимирова?',
               "расскажи деревья разбора",
               "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?"]

# %% colab={"base_uri": "https://localhost:8080/"} id="niysV3awx8eZ" outputId="f1bd3723-d09e-484a-e7f2-9853a3266480"
embeddings_1 = model.encode(sentences_1)['dense_vecs']
embeddings_2 = model.encode(sentences_2)['dense_vecs']

# %% id="zMsCd9Ri5jvL"
em_1 = mmm.encode(sentences_1)
em_2 = mmm.encode(sentences_2)

# %% colab={"base_uri": "https://localhost:8080/"} id="eOxM8uv74kOx" outputId="12966169-1e9a-4a76-8621-bacc66472348"
for i in range(len(embeddings_1)):
    for j in range(len(embeddings_2)):
        e_1 = embeddings_1[i]
        e_2 = embeddings_2[j]
        st_1 = em_1[i]
        st_2 = em_2[j]
        s_1 = sentences_1[i]
        s_2 = sentences_2[j]
        print(f'{s_1} AND {s_2}: ', e_1 @ e_2, st_1 @ st_2)

# %% id="cMfWelM72i4D"
