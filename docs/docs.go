// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag/v2"

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
        "/free": {
            "post": {
                "tags": [
                    "scooters"
                ],
                "summary": "Free the given scooter.",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "ClientID",
                        "name": "client_id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Scooter to free information",
                        "name": "Payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.FreePost"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    }
                }
            }
        },
        "/rent": {
            "post": {
                "tags": [
                    "scooters"
                ],
                "summary": "Rents the chosen scooter in given city.",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "ClientID",
                        "name": "client_id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Rental information details",
                        "name": "Payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.RentPost"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    }
                }
            }
        },
        "/scooters": {
            "get": {
                "tags": [
                    "scooters"
                ],
                "summary": "Gets scooters in the queried area of given city.",
                "parameters": [
                    {
                        "maxLength": 36,
                        "minLength": 36,
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "ClientID",
                        "name": "client_id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Ottawa",
                        "description": "City",
                        "name": "city",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 73.4,
                        "description": "Longitude of the center of the rectangle",
                        "name": "longitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 45.4,
                        "description": "Latitude of the center of the rectangle",
                        "name": "latitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 20000,
                        "description": "Height of the rectangle in meters",
                        "name": "height",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 25000,
                        "description": "Width of the rectangle in meters",
                        "name": "width",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.ScooterGet"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.FreePost": {
            "type": "object",
            "required": [
                "UUID"
            ],
            "properties": {
                "UUID": {
                    "type": "string"
                }
            }
        },
        "github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.RentPost": {
            "type": "object",
            "required": [
                "UUID",
                "city",
                "latitude",
                "longitude"
            ],
            "properties": {
                "UUID": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "github_com_NordSecurity-Interviews_BE-PatrykPasterny_internal_model.ScooterGet": {
            "type": "object",
            "properties": {
                "UUID": {
                    "type": "string"
                },
                "availability": {
                    "type": "boolean"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "model.ApiError": {
            "type": "object",
            "properties": {
                "Message": {
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
