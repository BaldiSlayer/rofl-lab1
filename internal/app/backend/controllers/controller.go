package controllers

type Controller struct {
	TRSParserClient interface {
		Parse(trs string) (string, error)
	}

	ModelClient interface {
		SendRequest(request string) (string, error)
	}
}
