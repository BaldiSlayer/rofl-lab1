package vdatabase

// CMem - реализация векторной базы данных с помощью chromem-go
type Mock struct {
}

func (cm *Mock) Init(baseURL string, contextFile string) error {
	return nil
}

func (cm *Mock) GetSimilar(question string) ([]string, error) {
	return []string{
		question,
	}, nil
}
