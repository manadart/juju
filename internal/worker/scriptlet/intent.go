// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package scriptlet

import "sync"

// IntentCollector accumulates intents declared by a scriptlet during
// execution. Nothing is applied until the scriptlet completes
// successfully. The collector is safe for concurrent use (Starform
// may invoke methods from multiple goroutines).
type IntentCollector struct {
	mu      sync.Mutex
	intents []Intent
}

// NewIntentCollector returns a new, empty intent collector.
func NewIntentCollector() *IntentCollector {
	return &IntentCollector{}
}

// Len returns the number of collected intents.
func (c *IntentCollector) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.intents)
}

// Intents returns a copy of all collected intents.
func (c *IntentCollector) Intents() []Intent {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]Intent, len(c.intents))
	copy(result, c.intents)
	return result
}

// StatusSet records an intent to set the application status.
func (c *IntentCollector) StatusSet(status, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.intents = append(c.intents, Intent{
		Type: IntentStatusSet,
		Data: StatusSetIntent{
			Status:  status,
			Message: message,
		},
	})
}

// StateSet records an intent to set a charm state key-value pair.
func (c *IntentCollector) StateSet(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.intents = append(c.intents, Intent{
		Type: IntentStateSet,
		Data: StateSetIntent{
			Key:   key,
			Value: value,
		},
	})
}

// StateDelete records an intent to delete a charm state key.
func (c *IntentCollector) StateDelete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.intents = append(c.intents, Intent{
		Type: IntentStateDelete,
		Data: StateDeleteIntent{
			Key: key,
		},
	})
}

// RelationSet records an intent to set relation data on an endpoint.
func (c *IntentCollector) RelationSet(endpoint string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Copy the map to prevent caller mutation.
	dataCopy := make(map[string]string, len(data))
	for k, v := range data {
		dataCopy[k] = v
	}
	c.intents = append(c.intents, Intent{
		Type: IntentRelationSet,
		Data: RelationSetIntent{
			Endpoint: endpoint,
			Data:     dataCopy,
		},
	})
}

// OpenPort records an intent to open a port on an endpoint.
func (c *IntentCollector) OpenPort(endpoint, protocol string, port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.intents = append(c.intents, Intent{
		Type: IntentOpenPort,
		Data: OpenPortIntent{
			Endpoint: endpoint,
			Protocol: protocol,
			Port:     port,
		},
	})
}

// ClosePort records an intent to close a port on an endpoint.
func (c *IntentCollector) ClosePort(endpoint, protocol string, port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.intents = append(c.intents, Intent{
		Type: IntentClosePort,
		Data: ClosePortIntent{
			Endpoint: endpoint,
			Protocol: protocol,
			Port:     port,
		},
	})
}

// Intent represents a single declared intent from a scriptlet.
type Intent struct {
	// Type identifies the kind of intent.
	Type IntentType

	// Data holds the intent-specific payload.
	Data interface{}
}

// IntentType identifies a specific kind of intent.
type IntentType string

const (
	IntentStatusSet  IntentType = "status-set"
	IntentStateSet   IntentType = "state-set"
	IntentStateDelete IntentType = "state-delete"
	IntentRelationSet IntentType = "relation-set"
	IntentOpenPort   IntentType = "open-port"
	IntentClosePort  IntentType = "close-port"
)

// StatusSetIntent holds data for a status-set intent.
type StatusSetIntent struct {
	Status  string
	Message string
}

// StateSetIntent holds data for a state-set intent.
type StateSetIntent struct {
	Key   string
	Value string
}

// StateDeleteIntent holds data for a state-delete intent.
type StateDeleteIntent struct {
	Key string
}

// RelationSetIntent holds data for a relation-set intent.
type RelationSetIntent struct {
	Endpoint string
	Data     map[string]string
}

// OpenPortIntent holds data for an open-port intent.
type OpenPortIntent struct {
	Endpoint string
	Protocol string
	Port     int
}

// ClosePortIntent holds data for a close-port intent.
type ClosePortIntent struct {
	Endpoint string
	Protocol string
	Port     int
}
