package mclient

type MockModelCient struct{}

func (mc *MockModelCient) Ask(request string) (string, error) {
	return request, nil
}
