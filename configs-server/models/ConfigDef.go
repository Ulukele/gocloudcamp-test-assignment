package models

type ConfigDef struct {
	Service string            `json:"service"`
	Data    map[string]string `json:"data"`
}
