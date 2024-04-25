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
        "/tax/calculations": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "parameters": [
                    {
                        "description": " ",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tax.taxCalculateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/tax.taxCalculateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError-array_validator_Field"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError-string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.ResponseError-array_validator_Field": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/validator.Field"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "http.ResponseError-string": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "errors": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "tax.allowance": {
            "type": "object",
            "properties": {
                "allowanceType": {
                    "type": "string"
                },
                "amount": {
                    "type": "number"
                }
            }
        },
        "tax.taxCalculateRequest": {
            "type": "object",
            "properties": {
                "allowances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.allowance"
                    }
                },
                "totalIncome": {
                    "type": "number",
                    "minimum": 0,
                    "example": 500000
                },
                "wht": {
                    "type": "number"
                }
            }
        },
        "tax.taxCalculateResponse": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number"
                }
            }
        },
        "validator.Field": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
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
	Title:            "Assessment Tax API",
	Description:      "This is a assessment tax api.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
