package dto

import (
	"encoding/json"
	"time"
)

type PipelineBundle struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Pipelines  []Pipeline `json:"pipelines"`
}

type Pipeline struct {
	Name               string      `json:"name"`
	Source             Component   `json:"source"`
	Queue              Component   `json:"queue"`
	OntologySinks      []Component `json:"ontologySinks"`
	SourceInterceptors []Component `json:"sourceInterceptors"`
	SinkInterceptors   []Component `json:"sinkInterceptors"`
}

type Component map[string]any

type CreatePipelineRuleRequest struct {
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description"`
	SourceType  string          `json:"source_type" binding:"required,oneof=file mysql api"`
	Status      string          `json:"status" binding:"omitempty,oneof=enabled disabled"`
	Config      json.RawMessage `json:"config" binding:"required" swaggertype:"object"`
}

type UpdatePipelineRuleRequest struct {
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description"`
	Status      string          `json:"status" binding:"required,oneof=enabled disabled"`
	Config      json.RawMessage `json:"config" binding:"required" swaggertype:"object"`
}

type PipelineRuleResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SourceType  string          `json:"source_type"`
	Status      string          `json:"status"`
	Version     int             `json:"version"`
	Config      json.RawMessage `json:"config" swaggertype:"object"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateServiceConnectionRequest struct {
	Name        string            `json:"name" binding:"required,max=255"`
	Description string            `json:"description"`
	SourceType  string            `json:"source_type" binding:"required,oneof=mysql api"`
	Status      string            `json:"status" binding:"omitempty,oneof=enabled disabled"`
	Config      json.RawMessage   `json:"config" binding:"required" swaggertype:"object"`
	Credentials map[string]string `json:"credentials"`
}

type UpdateServiceConnectionRequest struct {
	Name        string             `json:"name" binding:"required,max=255"`
	Description string             `json:"description"`
	Status      string             `json:"status" binding:"required,oneof=enabled disabled"`
	Config      json.RawMessage    `json:"config" binding:"required" swaggertype:"object"`
	Credentials *map[string]string `json:"credentials"`
}

type ServiceConnectionResponse struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	SourceType            string          `json:"source_type"`
	Status                string          `json:"status"`
	Config                json.RawMessage `json:"config" swaggertype:"object"`
	CredentialFields      []string        `json:"credential_fields"`
	CredentialsConfigured bool            `json:"credentials_configured"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}
