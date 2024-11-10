from faissDB import *

# new_questions = [
#     {
#         "question": "Является ли регулярным язык Дика с единственным типом скобок? (язык Дика - множество правильных скобочных структур вместе с пустой структурой, образующее язык над алфавитом {a,b}.)",
#         "answer": "Язык Дика не является регулярным. Докажем с помощью леммы о накачке. Предположим, что он регулярный, тогда по лемме о накачке существует n, являющееся длиной накачки. Возьмём последовательность из n открывающих, а затем n закрывающих скобок. Для неё существуют соответствующие x, y, z из леммы о накачке. Но так как |xy| <= n, то y состоит только из открывающих скобок, причём по условию леммы y не пустая. А значит при i = 2 в строке xy^iz получится больше открывающих скобок, чем закрывающих, то есть это будет не правильной скобочной последовательностью. Получили противоречие. Следовательно язык Дика с единственным типом скобок не является регулярным."},
#     {"question": "Что представляет собой язык, распознаваемый недетерминированным конечным автоматом (НКА)?",
#      "answer": "Язык, распознаваемый недетерминированным конечным автоматом (НКА) – это все такие слова, по которым существует хотя бы один путь из стартовой вершины в терминальную."},
#     {
#         "question": "Дан регулярный язык, опиши алгоритм нахождения кратчайшего слова, принадлежащего этому регулярному языку",
#         "answer": "Регулярный язык может быть задан с помощью конечного автомата. Так как автомат конечен, то мы можем его обойти (пройти через все состояния) за конечное время. Так как нам нужно найти самое короткое слово, то эта задача сводится к тому, что нам необходимо найти кратчайший путь от стартовой вершины до какой-либо терминальной. По определению это можно сделать с помощью поиска в ширину (bfs, breadth-first search). Поиск в ширину - алгоритм, находящий все кратчайшие пути от заданной вершины в невзвешенном графе. Запускаем поиск в ширину и выйдем из него, когда пришли в терминальное состояние. Так как нам нужно явно найти кратчайшее слово, то после этого делаем восстановление ответа. Это можно сделать используя дополнительную структуру данных (например массив prev), для того, чтобы для каждого состояния v хранить состояние u, из которого мы в него пришли. Мы можем пройти от найденной вершины по массиву prev, пока не придем в начальное состояние. Записав все переходы мы получим кратчайшее слово, принадлежащее регулярному языку."},
#     {"question": "Опиши алгоритм подсчета количества слов определенной длины в заданном регулярном языке",
#      "answer": "Обозначим регулярный язык за L и пусть длина слов, количество которых мы хотим найти - l. Так как язык L регулярен, то мы можем построить соответствующий ему конечный автомат A. Решим задачу с помощью динамического программирования. Пусть a_(q,i) – количество слов длины i, переводящих автомат A из начального состояния q0 в состояние q. Чтобы пересчитать эту величину, нужно просуммировать значения динамического программирования из предыдущего по длине слоя для всех состояний, из которых есть ребро в состояние q. Ответом является сумма элементов столбца, отвечающего за длину l, соответствующих терминальным вершинам."},
#     {"question": "Какими являются языки недетерминированных автоматов с магазинной памятью?",
#      "answer": "Языки недетерминированных автоматов с магазинной памятью являются контекстно-свободными. То есть эти языки могут быть заданы с помощью контекстно свободных грамматик."},
#     {"question": "Какая структура данных может описать магазинную память у автомата с магазинной памятью?",
#      "answer": "Магазинная память у автомата с магазинной памятью является стеком."},
#     {"question": "Опиши алгоритм нахождения эпсилон замыкания для каждой из вершин автомата?",
#      "answer": "ε-замыкание состояния q – это множество состояний, достижимых из q только по ε-переходам. Соотсветственно эпсилон замыкание для каждой из вершин автомата можно предподсчитать с помощью поиска в глубину (dfs) для каждой вершины."},
#     {"question": "Дай определение произведения двух автоматов",
#      "answer": "Прямым произведением двух ДКА A1=⟨Σ1,Q1,s1,T1,δ1⟩ и A2=⟨Σ2,Q2,s2,T2,δ2⟩ называется ДКА A=⟨Σ,Q,s,T,δ⟩, где: 1) Σ = Σ1∪Σ2, то есть он работает над пересечением алфавитов двух данных автоматов 2) Q = Q1×Q2, множество пар состояниий включает в себя состояния обоих автоматов 3) s =⟨s1,s2), стартуем с символов в обоих автоматах 4) T=T1×T2, терминальные состояния включают в себя терминальные состояния обоих автоматов 5) δ(⟨q1,q2⟩,c)=⟨δ1(q1,c),δ2(q2,c)⟩, то есть переходим по символу в обоих автоматах"},
#     {"question": "Как из недетерминированного конечного автомата A сделать pushdown automat B?",
#      "answer": "Для этого нужно заменить переход из состояния q в состояние p по символу x на такой же переход, только добавить z_0/z_0, где z_0 это дно стека. Из этого следует, что регулярные языки являются подмножеством МП-автоматных языков (языков автоматов с магазинной памятью)."},
#     {"question": "Теорема Клини",
#      "answer": "Теорема Клини гласит о том, что множество языков, принимаемых детерминированным конечным автоматом совпадает с множеством языком, принимаемых академическим регулярным выражением."},
#     {"question": "Какой язык называется префиксным (беспрефиксным)",
#      "answer": "Язык L называется префиксным, если для любого w не равного v из L не верно, что w – префикс v. Также такие языки называют беспрефиксными."},
#     {
#         "question": "Что можно сказать о языке L, который принимается детерминированным конечным автоматом с магазинной памятью по пустому стеку",
#         "answer": "Это значит, что язык L принимается детерминированным конечным автоматом с магазинной памятью по терминальному состоянию, а также язык L – префиксный"},
#     {
#         "question": "Является ли регулярным язык Дика с единственным типом скобок? (язык Дика - множество правильных скобочных структур вместе с пустой структурой, образующее язык над алфавитом {a,b}.)",
#         "answer": "Язык Дика не является регулярным. Докажем с помощью леммы о накачке. Предположим, что он регулярный, тогда по лемме о накачке существует n, являющееся длиной накачки. Возьмём последовательность из n открывающих, а затем n закрывающих скобок. Для неё существуют соответствующие x, y, z из леммы о накачке. Но так как |xy| <= n, то y состоит только из открывающих скобок, причём по условию леммы y не пустая. А значит при i = 2 в строке xy^iz получится больше открывающих скобок, чем закрывающих, то есть это будет не правильной скобочной последовательностью. Получили противоречие. Следовательно язык Дика с единственным типом скобок не является регулярным."},
#     {"question": "Является ли контекстно-свободным языком разность контекстно-свободного и регулярного языка?",
#      "answer": "Да, разность контекстно-свободного и регулярного языка является контекстно-свободным языком."},
#     {"question": "Что такое переход машины Тьюринга?",
#      "answer": "Переход машины Тьюринга — это функция, зависящая от состояния конечного управления и обозреваемого символа. За один переход машина Тьюринга должна выполнить следующие действия: изменить состояние, записать ленточный символ в обозреваемую клетку, сдвинуть головку влево или вправо."},
#     {"question": "Какие языки допускаются при помощи машины Тьюринга?",
#      "answer": "Языки, допустимые с помощью машины Тьюринга, называются рекурсивно перечислимыми, или РП-языками."},
#     {"question": "Опишите прием «память в состоянии» машины Тьюринга.",
#      "answer": "Память в состоянии — конечное управление можно использовать не только для представления позиции в «Программе» машины Тьюринга, но и для хранения конечного объема данных."},
#     {"question": "Опишите прием «Подпрограммы» машины Тьюринга.",
#      "answer": "Подпрограмма машины Тьюринга представляет собой множество состояний, выполняющее некоторый полезный процесс. Это множество включает в себя стартовое состояние и еще одно состояние, которое не имеет переходов и служит состоянием «возврата» для передачи управления какому-либо множеству состояний, вызвавшему данную подпрограмму. «Вызов» подпрограммы возникает везде, где есть переход в ее начальное состояние."},
#     {"question": "Можно ли запомнить позицию ленточной головки в позиции управления у машины Тьюринга?",
#      "answer": "Хотя позиции конечны в каждый момент времени, всё множество позиций может быть и бесконечным. Если состояние должно представлять любую позицию головки, то в состоянии должен быть компонент данных, имеющий любое целое в качестве значения. Из-за этого компонента множество состояний должно быть бесконечным, даже если только конечное число состояний используется в любой конечный момент времени. Определение же машин Тьюринга требует, чтобы множество состояний было конечным. Таким образом, запомнить позицию ленточной головки в конечном управлении нельзя."},
#     {"question": "Что такое счетчиковая машина?",
#      "answer": "Счетчиковые машины — это класс машин, обладающий возможностью запоминать конечное число целых чисел (счетчиков) и совершать различные переходы в зависимости от того, какие из счетчиков равны 0 (если таковые вообще есть). Счетчиковая машина может только прибавить 1 к счетчику или вычесть 1 из него, но отличить значения двух различных ненулевых счетчиков она не способна."},
#     {"question": "Какой язык допускается счетчиковой машиной?",
#      "answer": "Каждый язык, допускаемый счетчиковой машиной, рекурсивно перечислим. Причина в том, что счетчиковые машины являются частным случаем магазинных, а магазинные — частным случаем многоленточных машин Тьюринга, которые по теореме допускают только рекурсивно перечислимые языки."},
#     {"question": "Допускается ли любой рекурсивно перечислимый язык двухсчетчиковой машиной?",
#      "answer": "Для имитации машины Тьюринга и, следовательно, для допускания любого рекурсивно перечислимого языка, достаточно двух счетчиков. Для обоснования этого утверждения вначале доказывается, что достаточно трех счетчиков, а затем три счетчика имитируются с помощью двух."},
#     {"question": "Что такое универсальная машина Тьюринга?",
#      "answer": "Универсальной машиной Тьюринга называют машину Тьюринга, которая может заменить собой любую машину Тьюринга. Получив на вход программу и входные данные, она вычисляет ответ, который вычислила бы по входным данным машина Тьюринга, чья программа была дана на вход."},
#     {"question": "Какое время необходимо многоленточной машине Тьюринга для имитации шагов компьютера?",
#      "answer": "Рассмотрим компьютер, обладающий следующими свойствами: у него есть только инструкции, увеличивающие максимальную длину слова не более чем на один; у него есть только инструкции, которые многоленточная машина Тьюринга может выполнить на словах длиной k за O(k^2) или меньшее число шагов. Шаг — это выполнение одной инструкции. Таким образом, выполнение n шагов работы компьютера можно проимитировать на многоленточной машине Тьюринга с использованием не более O(n^3) шагов."},
#     {
#         "question": "Как связаны мощности следующих машин Тьюринга: многодорожечная машина Тьюринга, машина Тьюринга с односторонней лентой, многоленточная машина Тьюринга, недетерминированная машина Тьюринга?",
#         "answer": "Многодорожечная машина Тьюринга, машина Тьюринга с односторонней лентой, многоленточная машина Тьюринга и недетерминированная машина Тьюринга, несмотря на различия в их конструкции или правилах работы, обладают одинаковой вычислительной мощностью, то есть способны вычислить одни и те же классы функций. Различия между видами машин Тьюринга (например, между машинами с одним или несколькими лентами) могут повлиять на эффективность вычислений (время или пространство), но не на саму вычислительную мощность."},
#     {
#         "question": "Если проблема P1 неразрешима и ее можно свести к проблеме P2, то является ли проблема P2 неразрешимой?",
#         "answer": "Если проблему P1 можно свести к проблеме P2 и если P1 неразрешима, то и P2 неразрешима."},
#     {"question": "Что такое рандомизированная машина Тьюринга?",
#      "answer": "Рандомизированная машина Тьюринга — это вариант многоленточной машины Тьюринга. Первая лента, как обычно для многоленточных машин, содержит вход. Вторая лента также начинается непустыми клетками. В принципе, вся она содержит символы 0 и 1, выбранные с вероятностью 1/2. Вторая лента называется случайной лентой. Третья и последующие, если используются, вначале пусты и при необходимости выступают как рабочие."},
#     {"question": "Рекурсивные языки.",
#      "answer": "Языки, допускаемые машинами Тьюринга, называются рекурсивно-перечислимыми (РП), а РП-языки, допускаемые МТ, которые всегда останавливаются, — рекурсивными. “Разрешимость” есть синоним “рекурсивности”, однако языки чаще называются “рекурсивными”, а проблемы (которые представляют собой языки, интерпретируемые как вопросы) — “разрешимыми”. Если язык не является рекурсивным, то проблема, которую выражает этот язык, называется “неразрешимой”. Рекурсивный язык позволяет построить разрешающую функцию: т.е. МТ, возвращающую один из двух результатов (да-нет), и корректно завершающую работу."},
#     {
#         "question": "Рекурсивно-перечислимые языки. Примеры языков, которые являются рекурсивно-перечислимыми, но не рекурсивными.",
#         "answer": "Язык L является рекурсивно-перечислимым (РП-языком), если L = L(M) для некоторой машины Тьюринга M. Проблема останова машины Тьюринга является РП, но не рекурсивной. В действительности, определенная А. М. Тьюрингом машина допускала, не попадая в допускающее состояние, а останавливаясь. Для МТ M можно определить H(M) как множество входов w, на которых M останавливается независимо от того, допускает ли M вход w. Тогда проблема останова состоит в определении множества таких пар (M, w), у которых w принадлежит H(M). Это еще один пример проблемы/языка, которая является РП, но не рекурсивной."},
#     {"question": "Что такое язык диагонализации L_d",
#      "answer": "Язык диагонализации L_d — это множество всех цепочек w_i, не принадлежащих L(M_i). Понятие M_i, “i-й машины Тьюринга”. Это машина Тьюринга M, кодом которой является i-я двоичная цепочка w_i. В язык L_d входит каждая цепочка в алфавите {0, 1}, которая, будучи проинтерпретированной как код МТ, не принадлежит языку этой МТ. Язык L_d является хорошим примером не РП-языка, т.е. его не допускает ни одна машина Тьюринга."}
# ]
#
# add_new_questions(new_questions, filename="vectorized_data")
question = [
    {"question": "является ли язык Дика (язык правильных скобочных последовательностей) регулярным?", "answer": ""},
    {"question": "машина Тьюринга?", "answer": ""}
]
results = process_questions(question, use_saved=True)

for result in results:
    print(result)
    print("-" * 50)