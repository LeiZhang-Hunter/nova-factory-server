package convert

import (
	"fmt"
	"log"
	"nova-factory-server/app/utils/gin_mcp/pkg/types"
)

// generateInputSchemaExtension creates the JSON schema for the tool's input parameters.
// This is a simplified version using basic reflection and not an external library.
func generateInputSchemaExtension(route RouteProperties, registeredSchemas map[string]types.RegisteredSchemaInfo) *types.JSONSchema {
	// Base schema structure
	schema := &types.JSONSchema{
		Type:       "object",
		Properties: make(map[string]*types.JSONSchema),
		Required:   make([]string, 0),
	}
	properties := schema.Properties
	required := schema.Required

	// 1. Extract Path Parameters
	matches := PathParamRegex.FindAllStringSubmatch(route.Path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			paramName := match[1]
			properties[paramName] = &types.JSONSchema{
				Type:        "string",
				Description: fmt.Sprintf("Path parameter: %s", paramName),
			}
			required = append(required, paramName) // Path params are always required
		}
	}

	// 2. Incorporate Registered Query and Body Types
	schemaKey := route.HTTPMethod + " " + route.Path
	if schemaInfo, exists := registeredSchemas[schemaKey]; exists {
		if isDebugMode() {
			log.Printf("Using registered schema for %s", schemaKey)
		}

		// Reflect Query Parameters (if applicable for method and type exists)
		if (route.HTTPMethod == "GET" || route.HTTPMethod == "DELETE") && schemaInfo.QueryType != nil {
			reflectAndAddProperties(schemaInfo.QueryType, properties, &required, "query")
		}

		// Reflect Body Parameters (if applicable for method and type exists)
		if (route.HTTPMethod == "POST" || route.HTTPMethod == "PUT" || route.HTTPMethod == "PATCH") && schemaInfo.BodyType != nil {
			reflectAndAddProperties(schemaInfo.BodyType, properties, &required, "body")
		}
	}

	// Update the required slice in the main schema
	schema.Required = required

	// If no properties were added (beyond path params), handle appropriately.
	// Depending on the spec, an empty properties object might be required.
	// if len(properties) == 0 { // Keep properties map even if empty
	// 	// Return schema with empty properties
	// 	return schema
	// }

	return schema
}
