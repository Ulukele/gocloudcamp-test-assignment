package models

type Config struct {
	Service string            `json:"service"`
	Version string            `json:"version"`
	Data    map[string]string `json:"data"`
}
