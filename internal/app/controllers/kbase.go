package controllers

func (controller *Controller) knowledgeBase(question string) (string, error) {
	answer, err := controller.ModelClient.Ask(question)
	if err != nil {
		return "", err
	}

	return answer, nil
}
