package models

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Params  string `json:"params"`
}
