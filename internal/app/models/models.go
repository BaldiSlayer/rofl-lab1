package models

type QAPair struct {
	Question string `json:"question" yaml:"question"`
	Answer   string `json:"answer" yaml:"answer"`
}

type ModelRequest struct {
	Model      string
	UseContext bool
}
