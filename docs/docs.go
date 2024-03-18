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
            "name": "Alex Chinaev",
            "url": "https://vk.com/l.chinaev",
            "email": "ax.chinaev@yandex.ru"
        },
        "license": {
            "name": "AS IS (NO WARRANTY)"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/actors": {
            "get": {
                "description": "Gets all actors with related films.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Actors"
                ],
                "summary": "Gets actors.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "actors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.ActorWithFilms"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Modify an actor by id and retrieves a new actor.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Actors"
                ],
                "summary": "Modify an actor.",
                "parameters": [
                    {
                        "description": "Actor to modify",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ActorWithoutFilms"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "actors": {
                                            "$ref": "#/definitions/domain.ActorWithoutFilms"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new actor with the provided data.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Actors"
                ],
                "summary": "Adds a new actor.",
                "parameters": [
                    {
                        "description": "actor to add",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ActorToAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/actors/{id}": {
            "delete": {
                "description": "Deletes an actor by id with all its relations with films.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Actors"
                ],
                "summary": "Deletes an actor.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Actor id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/auth/login": {
            "post": {
                "description": "create user session and put it into cookie",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "login user",
                "parameters": [
                    {
                        "description": "user credentials",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Credentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/auth/logout": {
            "post": {
                "description": "delete current session and nullify cookie",
                "tags": [
                    "Auth"
                ],
                "summary": "logout user",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "add new user to db and return it id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "register user",
                "parameters": [
                    {
                        "description": "user credentials",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Credentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/films": {
            "get": {
                "description": "Gets all films descending sorted by rating (by default). Only one sort can be applied at a time. If several are applied, the priority is as follows: title, releaseDate, rating (by default).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Gets films.",
                "parameters": [
                    {
                        "enum": [
                            "Asc",
                            "Desc"
                        ],
                        "type": "string",
                        "description": "Direction of title sort. Sorting wont be applied if param isnt specified.",
                        "name": "sortTitle",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "Asc",
                            "Desc"
                        ],
                        "type": "string",
                        "description": "Direction of release date sort. Sorting wont be applied if param isnt specified.",
                        "name": "sortReleaseDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "films": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.FilmWithoutActors"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Modify a film by id and retrieves a new film.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Modify a film.",
                "parameters": [
                    {
                        "description": "Film to modify",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.FilmWithoutActors"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "film": {
                                            "$ref": "#/definitions/domain.FilmWithoutActors"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new film with provided data.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Adds a new film.",
                "parameters": [
                    {
                        "description": "film to add",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.FilmToAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/films/search": {
            "get": {
                "description": "Searches films by parts of its titles and parts of films names.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Searches films",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The string to be searched for",
                        "name": "searchStr",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "body": {
                                    "type": "object",
                                    "properties": {
                                        "films": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/domain.FilmWithoutActors"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/films/{id}": {
            "delete": {
                "description": "Deletes a film by id with all its relations with actors.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Deletes a film.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Film id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "err": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.ActorToAdd": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string",
                    "format": "date"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "$ref": "#/definitions/domain.Sex"
                }
            }
        },
        "domain.ActorToFilmAdd": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "domain.ActorWithFilms": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string",
                    "format": "date"
                },
                "films": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.FilmWithoutActors"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "$ref": "#/definitions/domain.Sex"
                }
            }
        },
        "domain.ActorWithoutFilms": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string",
                    "format": "date"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "$ref": "#/definitions/domain.Sex"
                }
            }
        },
        "domain.Credentials": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "domain.FilmToAdd": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ActorToFilmAdd"
                    }
                },
                "description": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "releaseDate": {
                    "type": "string",
                    "format": "date"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "domain.FilmWithoutActors": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "number"
                },
                "releaseDate": {
                    "type": "string",
                    "format": "date"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "domain.Sex": {
            "type": "string",
            "enum": [
                "M"
            ],
            "x-enum-varnames": [
                "M"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "FilmLib API",
	Description:      "API of the FilmLib project",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
