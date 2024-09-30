package controllers

import (
	"encoding/json"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/trsclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/vdatabase"
	"log/slog"
	"net/http"
)

// Controller - тип, который объединяет все ручки
// служит для удобного хранения данных, общих для всех ручек
type Controller struct {
	TRSParserClient trsclient.TRSParserClient
	// ModelClient - клиент к LLM
	ModelClient mclient.ModelClient
	// VectorDatabase - векторная база данных, используется для поиска ближайших по значению записей
	VectorDatabase vdatabase.VectorDatabase
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

	err := json.NewEncoder(row.w).Encode(map[string]string{
		"error": row.errorText,
	})
	if err != nil {
		slog.Error(row.errorText, "error", err)
	}
}
