package trsclient

type Mock struct{}

func (mc *Mock) Parse(trs string) (string, error) {
	return trs, nil
}
