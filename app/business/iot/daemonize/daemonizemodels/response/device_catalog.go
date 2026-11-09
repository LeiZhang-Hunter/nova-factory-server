package response

import "encoding/json"

// Catalog is an immutable, complete, gateway-scoped configuration snapshot.
type Catalog struct {
	GatewayID string   `json:"gatewayId"`
	Version   string   `json:"version"`
	Devices   []Device `json:"devices"`
}

// Device carries configuration, with no permission to enable local collection.
type Device struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Number       string          `json:"number"`
	Protocol     string          `json:"protocol"`
	Config       json.RawMessage `json:"config"`
	TemplateID   string          `json:"templateId"`
	TemplateName string          `json:"templateName"`
	Points       []Point         `json:"points"`
}

// Point describes one protocol measurement.
type Point struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Register       int64           `json:"register"`
	FunctionCode   int8            `json:"functionCode"`
	Sort           string          `json:"sort"`
	Unit           string          `json:"unit"`
	Precision      int64           `json:"precision"`
	Expression     string          `json:"expression"`
	DataFormat     string          `json:"dataFormat"`
	ProtocolConfig json.RawMessage `json:"protocolConfig,omitempty"`
}
