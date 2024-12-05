package mclient

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getContextFromQASlice(t *testing.T) {
	type args struct {
		contextQASlice []models.QAPair
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without context",
			args: args{
				contextQASlice: []models.QAPair{},
			},
			want: `
Вы - полезный ассистент, имеющий доступ к базе знаний, которому поручено отвечать на вопросы по теории формального языка.
Ничего не выдумывай. Если вы в чем-то не уверен, просто скажи, что не знаешь.



ВСЕ ТВОИ ОТВЕТЫ ОЧЕНЬ ВАЖНЫ ДЛЯ МЕНЯ. ТЫ ЧАСТЬ МОЕЙ ЛАБОРАТОРНОЙ РАБОТЫ, ЧЕМ ЛУЧШЕ ТЫ БУДЕШЬ ОТВЕЧАТЬ, ТЕМ БОЛЕЕ ХОРОШУЮ
ОЦЕНКУ МНЕ ПОСТАВЯТ. ПОЖАЛУЙСТА БУДЬ ОЧЕНЬ АККУРАТЕН В ОТВЕТАХ.
`,
		},
		{
			name: "without context",
			args: args{
				contextQASlice: []models.QAPair{
					{
						Question: "Что такое ТФЯ?",
						Answer:   "ТФЯ - предмет, который я к сожалению не знаю ((",
					},
					{
						Question: "Что такое ДКА?",
						Answer:   "Это детерминированный конечный автомат",
					},
					{
						Question: "Что такое НКА?",
						Answer:   "НКА это недетерминированный конечный автомат",
					},
				},
			},
			want: `
Вы - полезный ассистент, имеющий доступ к базе знаний, которому поручено отвечать на вопросы по теории формального языка.
Ничего не выдумывай. Если вы в чем-то не уверен, просто скажи, что не знаешь.



Отвечай на вопрос, основываясь исключительно на результатах поиска в базе знаний.

Все, что находится между следующими "контекстными" XML-блоками, извлекается из базы знаний, а не из диалога с пользователем. Пункты списка упорядочены по релевантности, поэтому первый из них является наиболее релевантным.
При ответе отдавай предпочтение контексту.

<context>

	<qa>
		<question>Что такое ТФЯ?</question>
		<answer>ТФЯ - предмет, который я к сожалению не знаю ((</answer>
	</qa>

	<qa>
		<question>Что такое ДКА?</question>
		<answer>Это детерминированный конечный автомат</answer>
	</qa>

	<qa>
		<question>Что такое НКА?</question>
		<answer>НКА это недетерминированный конечный автомат</answer>
	</qa>

</context>

Не упоминай в своем ответе базу знаний, контекст или результаты поиска. ОПИРАЙСЯ НА БАЗУ ЗНАНИЙ.


ВСЕ ТВОИ ОТВЕТЫ ОЧЕНЬ ВАЖНЫ ДЛЯ МЕНЯ. ТЫ ЧАСТЬ МОЕЙ ЛАБОРАТОРНОЙ РАБОТЫ, ЧЕМ ЛУЧШЕ ТЫ БУДЕШЬ ОТВЕЧАТЬ, ТЕМ БОЛЕЕ ХОРОШУЮ
ОЦЕНКУ МНЕ ПОСТАВЯТ. ПОЖАЛУЙСТА БУДЬ ОЧЕНЬ АККУРАТЕН В ОТВЕТАХ.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := getContextFromQASlice(tt.args.contextQASlice); got != tt.want {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
