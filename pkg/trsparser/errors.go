package trsparser

func (e *ParseError) Error() string {
	return e.Summary
}

