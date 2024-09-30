package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func (controller *Controller) knowledgeBase(question string) (string, error) {
	similar, err := controller.VectorDatabase.GetSimilar(question)
	if err != nil {
		return "", err
	}

	answer, err := controller.ModelClient.AskWithContext(question, similar)
	if err != nil {
		return "", err
	}

	return answer, nil
}

// KnowledgeBase - ручка, отвечает на запросы к базе знаний
func (controller *Controller) KnowledgeBase(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type QuestionData struct {
		Question string `json:"question"`
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

	var data QuestionData
	if err := json.Unmarshal(body, &data); err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to unmarshall data",
		})

		return
	}

	parseResult, err := controller.knowledgeBase(data.Question)
	if err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to get response from model",
		})

		return
	}

	_, err = fmt.Fprintf(w, parseResult)
	if err != nil {
		ErrorHandler(errorRow{
			w:         w,
			code:      http.StatusBadRequest,
			err:       err,
			errorText: "failed to write answer",
		})

		return
	}
}
