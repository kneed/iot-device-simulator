package models

type Strategy struct {
	Model

	Internal int `json:"internal"`
	Duration int `json:"duration"`
}