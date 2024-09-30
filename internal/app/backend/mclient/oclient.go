package mclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"
	"io"
	"net/http"
	"strings"
	"text/template"
)

const (
	llmModel = "llama3"

	bufferSize = 4096
)

var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
You are a helpful assistant with access to a knowlege base, tasked with answering questions about the world and its history, people, places and other things.

Answer the question in a very concise manner. Use an unbiased and journalistic tone. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
Answer the question solely based on the provided search results from the knowledge base. If the search results from the knowledge base are not relevant to the question at hand, just say that you don't know. Don't make anything up.

Anything between the following 'context' XML blocks is retrieved from the knowledge base, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
    {{- if . -}}
    {{- range $context := .}}
    - Question: {{.Question}} Answer: {{.Answer}} {{end}}
    {{- end}}
</context>
{{- end -}}

Don't mention the knowledge base, context or search results in your answer.
`))

type Ollama struct {
	ollamaBaseURL string
}

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type Response struct {
	Response string `json:"response"`
}

func NewOllamaApiClient(url string) *Ollama {
	var oll Ollama

	oll.ollamaBaseURL = url + "/api"

	return &oll
}

func (mc *Ollama) generateResponse(prompt string) (string, error) {
	reqBody := Request{
		Model:  llmModel,
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(mc.ollamaBaseURL+"/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var responseString string

	buffer := make([]byte, bufferSize)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			var tempResponse Response

			if err := json.Unmarshal(buffer[:n], &tempResponse); err == nil {
				responseString += tempResponse.Response
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("error reading response: %v", err)
		}
	}

	return responseString, nil
}

func (mc *Ollama) Ask(question string) (string, error) {
	return mc.generateResponse(question)
}

func (mc *Ollama) generateResponseWithContext(prompt string, answerContext string) (string, error) {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type RequestWithContext struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	type ResponseWithContext struct {
		Message Message `json:"message"`
	}

	reqBody := RequestWithContext{
		Model: llmModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: answerContext,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(mc.ollamaBaseURL+"/chat",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var responseString string

	buffer := make([]byte, bufferSize)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			var tempResponse ResponseWithContext

			if err := json.Unmarshal(buffer[:n], &tempResponse); err == nil {
				responseString += tempResponse.Message.Content
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("error reading response: %v", err)
		}
	}

	return responseString, nil
}

func (mc *Ollama) AskWithContext(question string, answerContext []models.QAPair) (string, error) {
	sb := &strings.Builder{}

	err := systemPromptTpl.Execute(sb, answerContext)
	if err != nil {
		return "", fmt.Errorf("failed to paste context to template: %v", err)
	}

	reply, err := mc.generateResponseWithContext(question, sb.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate response with context: %v", err)
	}

	return reply, nil
}
