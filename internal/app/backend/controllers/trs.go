package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TRSCheck - ручка, которая отвечает за распознавание TRS
func (controller *Controller) TRSCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
