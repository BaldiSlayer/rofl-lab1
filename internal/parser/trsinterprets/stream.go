package interprets

import "github.com/BaldiSlayer/rofl-lab1/internal/parser/models"

type stream struct {
	channel chan models.Lexem
	ok bool
	val models.Lexem
}

func (s *stream) peek() models.Lexem {
	if !s.ok {
		var ok bool
		s.val, ok = <-s.channel
		if !ok {
			s.val = models.NewEofLexem()
		}

		s.ok = true
	}
	return s.val
}

func (s *stream) next() models.Lexem {
	if s.ok {
		s.ok = false
		return s.val
	}

	val, ok := <-s.channel
	if !ok {
		return models.NewEofLexem()
	}
	return val
}
