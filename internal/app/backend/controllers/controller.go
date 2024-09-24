package controllers

type Controller struct {
	TRSParserClient interface {
		Parse(trs string) (string, error)
	}

	ModelClient interface {
		// Ask отправляет запрос к модели и возвращает ответ и ошибку, если что-то пошло не так
		Ask(request string) (string, error)
	}
}
