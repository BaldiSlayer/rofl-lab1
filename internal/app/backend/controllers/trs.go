package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/xretry"
	"github.com/julienschmidt/httprouter"
)

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

	parseResult := ""

	// add delay
	err = xretry.Retry(func() error {
		answer, err := controller.ModelClient.Ask(data.Request)
		if err != nil {
			return err
		}

		trsParserAns, err := controller.TRSParserClient.Parse(answer)
		if err != nil {
			return err
		}

		parseResult = trsParserAns

		return nil
	}).Count(3)

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
