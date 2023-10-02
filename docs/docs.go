// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Chris Developer",
            "email": "chrisd3v3l0p3r@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/book/": {
            "post": {
                "description": "Get details of a book",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Get a book",
                "parameters": [
                    {
                        "description": "Add new book",
                        "name": "models.Book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
                        }
                    }
                }
            }
        },
        "/book/{title}": {
            "get": {
                "description": "Get details of a book",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Get a book",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Title of the book",
                        "name": "title",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Book": {
            "type": "object",
            "required": [
                "autore",
                "categoria",
                "copertina",
                "genere",
                "id_copertina",
                "prezzo",
                "quantita",
                "summary",
                "titolo"
            ],
            "properties": {
                "autore": {
                    "type": "string"
                },
                "categoria": {
                    "type": "string"
                },
                "copertina": {
                    "type": "string"
                },
                "genere": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "id_copertina": {
                    "type": "string"
                },
                "prezzo": {
                    "type": "number"
                },
                "quantita": {
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 1
                },
                "summary": {
                    "type": "string",
                    "maxLength": 512
                },
                "titolo": {
                    "type": "string",
                    "maxLength": 255
                }
            }
        },
        "responses.ResponseErrorJSON": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "192.168.3.8:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Book Store API",
	Description:      "This API manage a cart",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
