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
        "/book": {
            "get": {
                "description": "Get details for every book",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Get books",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of the pagination",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseDatabase"
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
        },
        "/book/": {
            "post": {
                "description": "Adds new book",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Adds a book",
                "parameters": [
                    {
                        "description": "Adds new book",
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
        "/book/autore/{id}": {
            "put": {
                "description": "Update name of the writer.\nThe name of the writer can't be more than 64 characters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update writer",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Writer",
                        "name": "writer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/categoria/{id}": {
            "put": {
                "description": "Update the category of a book\nThe types of categories are: Best Seller, New Releases, and Best Offers.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update category",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Category",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/copertina/{id}": {
            "put": {
                "description": "Update the type cover of a book.\nThere are two types of cover: Hard Cover or Flexible Cover.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update cover",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Cover",
                        "name": "cover",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/genere/{id}": {
            "put": {
                "description": "Update the genre of a book.\nHere a list of all genre supported: Action, Adventure, Business, Cookbooks, Drama, Detective, Fantasy, Fiction, History, Horror, Romance, Psychology, Science Fiction, Short Stories, Thriller.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update genre",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Genre",
                        "name": "genre",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/id-copertina/{id}": {
            "put": {
                "description": "Update the id cover of a book",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update id cover",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Id Cover",
                        "name": "idcover",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/prezzo/{id}": {
            "put": {
                "description": "Update the price of a book.\nThe type of the price is a float (e.g. 15.90).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update price",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Price",
                        "name": "price",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/quantita/{id}": {
            "put": {
                "description": "Update the quantity of a book\nYou can choose the quantity from 1 to 5",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update quantity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Quantity",
                        "name": "quantity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/summary/{id}": {
            "put": {
                "description": "Update the summary of a book.\nThe summary of the book cannot be less than 512 characters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update summary",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Summary",
                        "name": "summary",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/titolo/{id}": {
            "put": {
                "description": "Update title of a book.\nThe title cannot be more than 255 characters.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Update title",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Title",
                        "name": "title",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
        },
        "/book/{id}": {
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
                        "type": "integer",
                        "description": "Title of the book",
                        "name": "id",
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
            },
            "delete": {
                "description": "Delete a book with same ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "Delete a book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.ResponseErrorJSON"
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
                    "type": "string",
                    "maxLength": 64,
                    "example": "Ruben Fabrizi"
                },
                "categoria": {
                    "type": "string",
                    "example": "New Releases"
                },
                "copertina": {
                    "type": "string",
                    "example": "Hard Cover"
                },
                "genere": {
                    "type": "string",
                    "example": "Romance"
                },
                "id": {
                    "type": "integer"
                },
                "id_copertina": {
                    "type": "integer",
                    "example": 564
                },
                "prezzo": {
                    "type": "number",
                    "example": 15.9
                },
                "quantita": {
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 1,
                    "example": 2
                },
                "summary": {
                    "type": "string",
                    "maxLength": 512,
                    "example": "Lorem Ipsum is simply dummy text of the printing and typesetting industry."
                },
                "titolo": {
                    "type": "string",
                    "maxLength": 255,
                    "example": "Il silenzio di un mare in tempesta"
                }
            }
        },
        "responses.PaginationDatabase": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                },
                "total_record": {
                    "type": "integer"
                }
            }
        },
        "responses.ResponseDatabase": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Book"
                    }
                },
                "paging": {
                    "description": "struct embedding",
                    "allOf": [
                        {
                            "$ref": "#/definitions/responses.PaginationDatabase"
                        }
                    ]
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
	Host:             "192.168.3.6:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Book Store API",
	Description:      "API for manages a book store",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
