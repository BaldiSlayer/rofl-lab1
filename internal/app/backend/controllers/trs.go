package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	fmt.Fprintf(w, "test")
}
