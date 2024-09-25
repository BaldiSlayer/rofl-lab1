package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/xretry"
	"github.com/julienschmidt/httprouter"
)

const (
	delayBetweenParseRetries = 3 * time.Second
	countParseRetries        = 3
)

// trsCheck проверяет завершается ли система переписывания термов.
// (Если только это, то надо будет поставить bool)
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

// TRSCheck - ручка, которая отвечает за распознавание TRS
func (controller *Controller) TRSCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type RequestData struct {
		Request string `json:"request"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to read request body",
		})

		return
	}
	defer r.Body.Close()

	var data RequestData
	if err := json.Unmarshal(body, &data); err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to unmarshall data",
		})

		return
	}

	parseResult, err := controller.trsCheck(data.Request)
	if err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to parse trs",
		})

		return
	}

	fmt.Fprintf(w, parseResult)
}
