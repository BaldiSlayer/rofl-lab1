package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/trsclient"
)

// Controller - тип, который объединяет все ручки
// служит для удобного хранения данных, общих для всех юзкейсов
type Controller struct {
	TRSParserClient trsclient.TRSParserClient
	// ModelClient - клиент к LLM
	ModelClient mclient.ModelClient
}
