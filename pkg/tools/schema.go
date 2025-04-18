package tools

import (
	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// GenerateSchema generates a JSON schema for a given type T
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
