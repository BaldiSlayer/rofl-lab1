package trsclient

type TRSParserClient interface {
	// Parse выполняет парсинг TRS, которая была выделена с помощью модели
	Parse(trs string) (string, error)
}
