package dto

type Log struct {
	Project string `json:"project"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
