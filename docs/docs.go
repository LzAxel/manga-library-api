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
        "/api/heartbeat": {
            "get": {
                "tags": [
                    "Metrics"
                ],
                "summary": "Heartbeat metric",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/manga": {
            "get": {
                "tags": [
                    "Manga"
                ],
                "summary": "Get Manga by ID or Slug",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Manga ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Manga slug",
                        "name": "slug",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Manga"
                        }
                    },
                    "400": {
                        "description": "Invalid input (One of params is required)"
                    },
                    "404": {
                        "description": "Manga not found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Manga"
                ],
                "summary": "Create Manga",
                "parameters": [
                    {
                        "description": "Add manga",
                        "name": "manga",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateMangaDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Created manga ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/manga/latest": {
            "get": {
                "tags": [
                    "Manga"
                ],
                "summary": "Get Latest Manga",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Manga"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/manga/{id}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Manga"
                ],
                "summary": "Delete Manga by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Manga ID",
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
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Manga"
                ],
                "summary": "Update Manga by ID",
                "parameters": [
                    {
                        "description": "Update manga",
                        "name": "manga",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateMangaDTO"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Manga ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "Sign in via username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authorization"
                ],
                "summary": "Sign In",
                "parameters": [
                    {
                        "description": "Auth Sign In Input",
                        "name": "auth",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.LoginUserDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "Sign up via username and password",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authorization"
                ],
                "summary": "Sign up",
                "parameters": [
                    {
                        "description": "Auth Sign Up Input",
                        "name": "auth",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "409": {
                        "description": "Username is already exists"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Chapter": {
            "type": "object",
            "properties": {
                "isPublished": {
                    "type": "boolean"
                },
                "number": {
                    "type": "integer"
                },
                "pageCount": {
                    "type": "integer"
                },
                "uploadedAt": {
                    "type": "string"
                },
                "uploaderId": {
                    "type": "string"
                },
                "volume": {
                    "type": "integer"
                }
            }
        },
        "domain.CreateMangaDTO": {
            "type": "object",
            "required": [
                "ageRating",
                "author",
                "description",
                "previewUrl",
                "releaseYear",
                "tags",
                "title"
            ],
            "properties": {
                "ageRating": {
                    "type": "integer"
                },
                "alternativeTitles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "author": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "previewUrl": {
                    "type": "string"
                },
                "releaseYear": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "domain.CreateUserDTO": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "domain.LoginUserDTO": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "domain.Manga": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "ageRating": {
                    "type": "integer"
                },
                "alternativeTitles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "author": {
                    "type": "string"
                },
                "chapters": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Chapter"
                    }
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "isPublished": {
                    "type": "boolean"
                },
                "likeCount": {
                    "type": "integer"
                },
                "previewUrl": {
                    "type": "string"
                },
                "releaseYear": {
                    "type": "integer"
                },
                "slug": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                },
                "uploaderId": {
                    "type": "string"
                }
            }
        },
        "domain.UpdateMangaDTO": {
            "type": "object",
            "properties": {
                "ageRating": {
                    "type": "integer"
                },
                "alternativeTitles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "author": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "isPublished": {
                    "type": "boolean"
                },
                "previewUrl": {
                    "type": "string"
                },
                "releaseYear": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Manga Library API",
	Description:      "API Server for Manga Library App",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
