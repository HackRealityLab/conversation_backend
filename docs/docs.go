// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/conversation/file": {
            "post": {
                "description": "Принимает аудиофайл разговора",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conversation"
                ],
                "summary": "Загрузка аудиофайла разговора",
                "operationId": "load-conversation-file",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        },
        "/conversation/file/send_ai/{name}": {
            "post": {
                "description": "Отправка аудиофайла нейронке",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conversation"
                ],
                "summary": "Отправка аудиофайла нейронке",
                "operationId": "send-file-ai",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        },
        "/conversation/file/{name}": {
            "get": {
                "description": "Получение аудиофайла разговора по его названию. Возвращает файл",
                "produces": [
                    "multipart/form-data"
                ],
                "tags": [
                    "conversation"
                ],
                "summary": "Получение аудиофайла разговора по его названию",
                "operationId": "get-conversation-file",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        },
        "/conversation/records": {
            "get": {
                "description": "Отправка записей",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "record"
                ],
                "summary": "Получение записей",
                "operationId": "get-records",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Record"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        },
        "/conversation/records/{id}": {
            "get": {
                "description": "Получение записи",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "record"
                ],
                "summary": "Получение записи",
                "operationId": "get-record",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Record"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        },
        "/conversation/text": {
            "post": {
                "description": "Принимает текст разговора",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conversation"
                ],
                "summary": "Загрузка текста разговора",
                "operationId": "load-conversation-text",
                "parameters": [
                    {
                        "description": "текст сообщения",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ConversationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Record": {
            "type": "object",
            "properties": {
                "audioName": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "goodPercent": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "isOk": {
                    "type": "boolean"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dto.ConversationRequest": {
            "type": "object",
            "required": [
                "text"
            ],
            "properties": {
                "text": {
                    "type": "string"
                }
            }
        },
        "response.Body": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Hackathon API",
	Description:      "API Server for Hackathon",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
