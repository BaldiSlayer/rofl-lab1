package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Controller - тип, который объединяет все ручки
// служит для удобного хранения данных, общих для всех ручек
type Controller struct {
	TRSParserClient interface {
		// Parse выполняет парсинг TRS, которая была выделена с помощью модели
		Parse(trs string) (string, error)
	}

	ModelClient interface {
		// Ask отправляет запрос к модели
		Ask(request string) (string, error)
	}
}

type errorRow struct {
	w    http.ResponseWriter
	code int

	err       error
	errorText string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler(row errorRow) {
	slog.Error(row.errorText, "error", row.err)

	row.w.WriteHeader(row.code)
	json.NewEncoder(row.w).Encode(map[string]string{
		"error": row.errorText,
	})
}
