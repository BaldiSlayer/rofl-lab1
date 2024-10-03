package lexer

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
)

type Lexer struct {
	Text  string
	Lexem []models.Lexem
}

func (p *Lexer) appendLex(index int, lexType models.LexemType, str string) {
	p.Lexem = append(p.Lexem, models.Lexem{
		//index: index,
		LexemType: lexType,
		Str:       str,
	})
}

func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (p *Lexer) Process() error {
	runes := []rune(p.Text)
	p.Lexem = make([]models.Lexem, 0, len(runes))

	lexVariables := []rune("variables")

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case ' ': // пробел и таб пропустить
			continue
		case '\t':
			continue
		case '-':
			p.appendLex(i, models.LexSEPARATOR, "-")
			for i < len(runes) && runes[i] == '-' {
				i++
			}
		case '=':
			p.appendLex(i, models.LexEQ, "=")
		case ',':
			p.appendLex(i, models.LexCOMMA, ",")
		case '+':
			p.appendLex(i, models.LexADD, "+")
		case '*':
			p.appendLex(i, models.LexMUL, "*")
		case '{':
			p.appendLex(i, models.LexLCB, "{")
		case '}':
			p.appendLex(i, models.LexRCB, "}")
		case '(':
			p.appendLex(i, models.LexLB, "(")
		case ')':
			p.appendLex(i, models.LexRB, ")")
		default:
			if runes[i] == '\n' || runes[i] == '\r' { // если перевод строки(причем могут быть 2), добавить лексему перевод строки
				p.appendLex(i, models.LexEOL, "\n")
				if i < len(runes)-1 && (runes[i] == '\n' && runes[i+1] == '\r' || runes[i] == '\r' && runes[i+1] == '\n') {
					i++
				}
			} else if isLetter(runes[i]) { // если встретилась буква
				if runes[i] == 'v' && i+len(lexVariables) < len(runes) { // проверяем на "variables"
					wordVariablesFound := true
					j := 0
					for ; j < len(lexVariables); j++ {
						if lexVariables[j] != runes[i+j] {
							wordVariablesFound = false
							break
						}
					}
					if wordVariablesFound { // если найдено слово, добавляем и пропускаем
						p.appendLex(i, models.LexVAR, "variables")
						i += len(lexVariables) - 1
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						p.appendLex(i, models.LexLETTER, string(runes[i]))
					}
				} else { // если найденная буква не v, то добавляем букву
					p.appendLex(i, models.LexLETTER, string(runes[i]))
				}
			} else if isDigit(runes[i]) {
				start_index := i
				for i+1 < len(runes) && isDigit(runes[i+1]) {
					i++
				}
				p.appendLex(i, models.LexNUM, string(runes[start_index:i+1]))
			} else {
				return fmt.Errorf("Неизвестный символ на позиции %d:%c", i, runes[i])
			}
		}
	}
	return nil
}
