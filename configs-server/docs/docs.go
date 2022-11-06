// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/config/": {
            "get": {
                "description": "Responds with the config as JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Reads the latest version of the config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search config by service",
                        "name": "service",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Config"
                        }
                    }
                }
            },
            "put": {
                "description": "Responds with the config as JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Creates new config version",
                "parameters": [
                    {
                        "description": "Config JSON",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfigDef"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Config"
                        }
                    }
                }
            },
            "post": {
                "description": "Responds with the config as JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Create config",
                "parameters": [
                    {
                        "description": "Config JSON",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfigDef"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Config"
                        }
                    }
                }
            },
            "delete": {
                "description": "Responds with the all configs versions as JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Deletes all config versions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "delete configs by service",
                        "name": "service",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Config"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Config": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "service": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "models.ConfigDef": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "service": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}