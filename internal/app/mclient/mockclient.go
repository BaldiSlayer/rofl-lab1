package mclient

type Mock struct{}

func (mc *Mock) Ask(question string) (string, error) {
	return question, nil
}

func (mc *Mock) AskWithContext(question string) (string, error) {
	return question, nil
}
