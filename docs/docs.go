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
        "/api/banknote": {
            "get": {
                "description": "Banknotes List",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Banknotes List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Query string to filter banknotes by nominal",
                        "name": "banknote_name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ds.BanknoteList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/banknotes": {
            "post": {
                "description": "Add a new banknote with image, nominal, currency",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Add new banknote",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Banknote image",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Banknote nominal",
                        "name": "nominal",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Banknote description",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Banknote currency",
                        "name": "currency",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a banknote with the given ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Delete banknote by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Banknote ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/banknotes/": {
            "put": {
                "description": "Updates a banknote with the given ID",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Update banknote by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "nominal",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "currency",
                        "name": "IIN",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "image",
                        "name": "image",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/banknotes/request": {
            "post": {
                "description": "Adds a banknote to a operation request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Add banknote to request",
                "parameters": [
                    {
                        "description": "Добавление банкноты",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.AddToBanknoteID"
                        }
                    },
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/banknotes/{id}": {
            "get": {
                "description": "Banknote By ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Banknotes"
                ],
                "summary": "Banknotes By ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Banknotes ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ds.Banknote"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/operation-request-banknote": {
            "put": {
                "description": "Update money Operation Banknote by client",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operation_Banknote"
                ],
                "summary": "Update money Operation Banknote",
                "parameters": [
                    {
                        "description": "Update quantity Operation Banknote",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.OperationBanknote"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "update",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Deletes a banknote from a request based on the user ID and banknote ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operation_Banknote"
                ],
                "summary": "Delete banknote from request",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "banknote ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/operations": {
            "get": {
                "description": "Retrieves a list of Operation requests based on the provided parameters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Get list of Operation requests",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Operation request status",
                        "name": "status_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date in the format '2006-01-02T15:04:05Z'",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date in the format '2006-01-02T15:04:05Z'",
                        "name": "end_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ds.Operation"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "put": {
                "description": "Update Operation by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Update Operation by admin",
                "parameters": [
                    {
                        "description": "updated Assembly",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.Operation"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Deletes a operation request for the given user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Delete operation request by user ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/operations/form": {
            "put": {
                "description": "Form Banknote by client",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Form Banknote by client",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Operation form ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/operations/{id}": {
            "get": {
                "description": "Retrieves a operation request with the given ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Get operation request by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Operation Request ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/user/logout": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Завершение сеанса текущего пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пользователи"
                ],
                "summary": "Выход пользователя",
                "responses": {
                    "200": {
                        "description": "Успешный выход",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    },
                    "401": {
                        "description": "Неверные учетные данные",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    }
                }
            }
        },
        "/api/user/signIn": {
            "post": {
                "description": "Вход нового пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пользователи"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Детали входа",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.RegisterReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация",
                        "schema": {
                            "$ref": "#/definitions/ds.LoginSwaggerResp"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    },
                    "401": {
                        "description": "Неверные учетные данные",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.errorResp"
                        }
                    }
                }
            }
        },
        "/api/user/signUp": {
            "post": {
                "description": "Регистрация нового пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пользователи"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Детали регистрации",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.RegisterReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/operations/updateStatus": {
            "put": {
                "description": "Updates the status of a operation request with the given ID on \"завершен\"/\"отклонен\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operation"
                ],
                "summary": "Update operation request status by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Request ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update status",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ds.NewStatus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "ds.AddToBanknoteID": {
            "type": "object",
            "properties": {
                "banknote_id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "ds.Banknote": {
            "type": "object",
            "properties": {
                "banknote_id": {
                    "type": "integer"
                },
                "currency": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "nominal": {
                    "type": "number"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "ds.BanknoteList": {
            "type": "object",
            "properties": {
                "banknotes_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ds.Banknote"
                    }
                },
                "draft_id": {
                    "type": "integer"
                }
            }
        },
        "ds.LoginSwaggerResp": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "string"
                },
                "token_type": {
                    "type": "string"
                }
            }
        },
        "ds.NewStatus": {
            "type": "object",
            "properties": {
                "operation_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "ds.Operation": {
            "type": "object",
            "properties": {
                "comletion_at": {
                    "type": "string"
                },
                "crated_at": {
                    "type": "string"
                },
                "formation_at": {
                    "type": "string"
                },
                "is_income": {
                    "type": "string"
                },
                "moderator": {
                    "$ref": "#/definitions/ds.User"
                },
                "moderator_id": {
                    "type": "integer"
                },
                "operation_banknote": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ds.OperationBanknote"
                    }
                },
                "operation_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "status_check": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/ds.User"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "ds.OperationBanknote": {
            "type": "object",
            "properties": {
                "banknote": {
                    "$ref": "#/definitions/ds.Banknote"
                },
                "banknote_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "operation_id": {
                    "type": "integer"
                },
                "opration": {
                    "$ref": "#/definitions/ds.Operation"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "ds.RegisterReq": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "ds.Role": {
            "type": "integer",
            "enum": [
                1,
                2,
                3
            ],
            "x-enum-varnames": [
                "Buyer",
                "Moderator",
                "Admin"
            ]
        },
        "ds.User": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                },
                "login": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "registrationDate": {
                    "type": "string"
                },
                "role": {
                    "$ref": "#/definitions/ds.Role"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "handler.errorResp": {
            "type": "object",
            "properties": {
                "error_description": {
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
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
