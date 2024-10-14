package controllers

import (
	"time"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/xretry"
)

const (
	delayBetweenParseRetries = 3 * time.Second
	countParseRetries        = 3
)

// trsCheck проверяет завершается ли система переписывания термов.
// (Если только это, то надо будет поставить bool)
// TODO: подкрутить интерфейс
func (controller *Controller) trsCheck(request string) (string, error) {
	parseResult := ""

	// пытаемся распарсить несколько раз
	err := xretry.Retry(func() error {
		answer, err := controller.ModelClient.Ask(request)
		if err != nil {
			return err
		}

		trsParserAns, err := controller.TRSParserClient.Parse(answer)
		if err != nil {
			return err
		}

		parseResult = trsParserAns

		return nil
	}).WithDelay(delayBetweenParseRetries).Count(countParseRetries)

	return parseResult, err
}
