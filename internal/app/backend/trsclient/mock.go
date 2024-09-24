package trsclient

type MockTRSClient struct{}

func (mc *MockTRSClient) Parse(trs string) (string, error) {
	return "", nil
}
