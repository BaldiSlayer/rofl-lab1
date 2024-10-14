package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/trsclient"
)

// Controller - тип, который объединяет все ручки
// служит для удобного хранения данных, общих для всех юзкейсов
type Controller struct {
	TRSParserClient trsclient.TRSParserClient
	// ModelClient - клиент к LLM
	ModelClient mclient.ModelClient
}
