package rest

import (
	"encoding/json"
	"errors"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"

	"github.com/getkin/kin-openapi/openapi3"
)

//router.GET("/swagger.json", func(c *gin.Context) {
//	endpoints := loadEndpointDescriptions() // Load or update your endpoint descriptions
//	swagger := generateSwaggerDocument(endpoints)
//	c.JSON(http.StatusOK, swagger)
//})

func GenerateSwaggerDocument(endpoints []*endpoint.Endpoint) (swagger *openapi3.T, err error) {
	swagger = &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "DTAC Agent Dynamic API",
			Version: "1.0.0",
		},
		Paths: &openapi3.Paths{},
	}

	for _, endpoint := range endpoints {
		path, err := convertToSwaggerPath(endpoint)
		if err != nil {
			return nil, err
		}
		swagger.Paths.Set(endpoint.Path, path)
	}

	return swagger, nil
}

// convertToSwaggerPath converts an Endpoint to a Swagger path object
func convertToSwaggerPath(endpoint *endpoint.Endpoint) (*openapi3.PathItem, error) {
	pathItem := &openapi3.PathItem{}

	// Create a new Operation object common for all actions
	operation := openapi3.NewOperation()
	operation.Summary = endpoint.Description
	operation.Description = endpoint.Description

	// Set up default responses (customize this as per your API's requirements)
	operation.Responses = openapi3.NewResponses(openapi3.WithName("200", openapi3.NewResponse().WithDescription("OK")))

	switch endpoint.Action {
	case "read":
		pathItem.Get = operation
	case "write":
		pathItem.Put = operation
	case "create":
		pathItem.Post = operation
	case "delete":
		pathItem.Delete = operation
	default:
		return nil, errors.New("unsupported action")
	}

	// Add description, parameters, request body, and responses
	operation.Summary = endpoint.Description
	operation.Description = endpoint.Description
	operation.Parameters = openapi3.Parameters{}
	operation.RequestBody = &openapi3.RequestBodyRef{}

	// Convert JSON Schema strings to OpenAPI schema objects
	if endpoint.ExpectedParametersSchema != "" {
		//schema, err := convertJSONSchemaToOpenAPI(endpoint.ExpectedParametersSchema)
		//if err != nil {
		//	return nil, err
		//}
		//// Add parameters schema to operation
		//// You may need to adjust this part based on how your schema defines parameters
	}

	if endpoint.ExpectedBodySchema != "" {
		schema, err := convertJSONSchemaToOpenAPI(endpoint.ExpectedBodySchema)
		if err != nil {
			return nil, err
		}
		// Add body schema to operation
		operation.RequestBody = &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().WithContent(openapi3.NewContentWithJSONSchema(schema)),
		}
	}

	// Add operation to the correct HTTP method in pathItem
	switch endpoint.Action {
	case "read":
		pathItem.Get = operation
	case "write":
		pathItem.Put = operation
	case "create":
		pathItem.Post = operation
	case "delete":
		pathItem.Delete = operation
	}

	return pathItem, nil
}

// convertJSONSchemaToOpenAPI converts a JSON Schema string to an OpenAPI schema object
func convertJSONSchemaToOpenAPI(jsonSchemaStr string) (*openapi3.Schema, error) {
	var schema openapi3.Schema
	err := json.Unmarshal([]byte(jsonSchemaStr), &schema)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}
