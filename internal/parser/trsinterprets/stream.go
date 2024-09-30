package interprets

import "github.com/BaldiSlayer/rofl-lab1/internal/trs"

type stream struct {
	channel chan trs.Lexem
	ok bool
	val trs.Lexem
}

func (s *stream) peek() trs.Lexem {
	if !s.ok {
		var ok bool
		s.val, ok = <-s.channel
		if !ok {
			s.val = trs.NewEofLexem()
		}

		s.ok = true
	}
	return s.val
}

func (s *stream) next() trs.Lexem {
	if s.ok {
		s.ok = false
		return s.val
	}

	val, ok := <-s.channel
	if !ok {
		return trs.NewEofLexem()
	}
	return val
}
