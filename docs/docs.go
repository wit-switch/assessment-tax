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
        "/admin/deductions/k-receipt": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "parameters": [
                    {
                        "description": " ",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.updateTaxDeductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/admin.updateKReceiptDeductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-array_validator_Field"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    }
                }
            }
        },
        "/admin/deductions/personal": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "parameters": [
                    {
                        "description": " ",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.updateTaxDeductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/admin.updatePersonalDeductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-array_validator_Field"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    }
                }
            }
        },
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
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-array_validator_Field"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    }
                }
            }
        },
        "/tax/calculations/upload-csv": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "parameters": [
                    {
                        "type": "file",
                        "description": " ",
                        "name": "taxFile",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/tax.texes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-array_validator_Field"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "admin.updateKReceiptDeductResponse": {
            "type": "object",
            "properties": {
                "kReceipt": {
                    "type": "number"
                }
            }
        },
        "admin.updatePersonalDeductResponse": {
            "type": "object",
            "properties": {
                "personalDeduction": {
                    "type": "number"
                }
            }
        },
        "admin.updateTaxDeductRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-array_validator_Field": {
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
        "github.com_wit-switch_assessment-tax_internal_handler_http.ResponseError-string": {
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
            "required": [
                "allowanceType"
            ],
            "properties": {
                "allowanceType": {
                    "type": "string",
                    "example": "donation"
                },
                "amount": {
                    "type": "number",
                    "minimum": 0,
                    "example": 200000
                }
            }
        },
        "tax.taxCSV": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number"
                },
                "taxRefund": {
                    "type": "number"
                },
                "totalIncome": {
                    "type": "number"
                }
            }
        },
        "tax.taxCalculateRequest": {
            "type": "object",
            "required": [
                "allowances"
            ],
            "properties": {
                "allowances": {
                    "type": "array",
                    "uniqueItems": true,
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
                    "type": "number",
                    "minimum": 0,
                    "example": 25000
                }
            }
        },
        "tax.taxCalculateResponse": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number"
                },
                "taxLevel": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.taxLevel"
                    }
                },
                "taxRefund": {
                    "type": "number"
                }
            }
        },
        "tax.taxLevel": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "string"
                },
                "tax": {
                    "type": "number"
                }
            }
        },
        "tax.texes": {
            "type": "object",
            "properties": {
                "texes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.taxCSV"
                    }
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
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
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
	Description:      "BasicAuth protects our entity endpoints.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
