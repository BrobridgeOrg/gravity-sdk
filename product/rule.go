package product

import "time"

// Rule defines the structure for a rule associated with a data product.
// It specifies how events are processed and managed within the system.
type Rule struct {
	ID            string                 `json:"id"`                // A unique identifier for the rule.
	Name          string                 `json:"name"`              // A human-readable name for the rule.
	Description   string                 `json:"desc"`              // A brief description of what the rule does.
	Event         string                 `json:"event"`             // The type of event that triggers this rule.
	Product       string                 `json:"product"`           // The name of the product associated with this rule.
	Method        string                 `json:"method"`            // The method or action to be performed when the rule is triggered.
	PrimaryKey    []string               `json:"primaryKey"`        // A list of fields used as the primary key in the rule processing.
	SchemaConfig  map[string]interface{} `json:"schema,omitempty"`  // Optional configuration for the data schema related to the rule.
	HandlerConfig *HandlerConfig         `json:"handler,omitempty"` // Optional configuration for the handler responsible for executing the rule.
	Enabled       bool                   `json:"enabled"`           // Flag indicating whether the rule is currently active or not.
	CreatedAt     time.Time              `json:"createdAt"`         // Timestamp indicating when the rule was created.
	UpdatedAt     time.Time              `json:"updatedAt"`         // Timestamp indicating the last update to the rule.
}

func NewRule() *Rule {
	return &Rule{}
}

type HandlerConfig struct {
	Type   string `json:"type"`
	Script string `json:"script"`
}
