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
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data RequestData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Welcome! Received request: %s\n", data.Request)
}
