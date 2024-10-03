package trsparser

func (e *ParseError) Error() string {
	return e.Summary
}

func toParseError(err error) error {
	// TODO: convert
	return err
}
