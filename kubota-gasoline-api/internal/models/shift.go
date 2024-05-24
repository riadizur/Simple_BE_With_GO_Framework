package models

type Shift struct {
	ID              int    `json:"id"`
	Date            string `json:"date"`
	Start           string `json:"start"`
	Finish          string `json:"finish,omitempty"`
	IsAlreadyFinish bool   `json:"isAlreadyFinish"`
}
