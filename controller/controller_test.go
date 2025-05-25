package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew tests the constructor function
func TestNew(t *testing.T) {
	// Arrange
	commit := "abc123"
	version := "1.0.0"
	delivery := NewMockDelivery(t)
	ingestion := NewMockIngestion(t)

	// Act
	ctrl := New(commit, version, delivery, ingestion)

	// Assert
	assert.NotNil(t, ctrl)
	assert.Equal(t, commit, ctrl.commit)
	assert.Equal(t, version, ctrl.version)
	assert.Equal(t, delivery, ctrl.delivery)
	assert.Equal(t, ingestion, ctrl.ingestion)
}

// MockDeliveryWithFields is used to expose fields for test assertions
type MockDeliveryWithFields struct {
	delivery  Delivery
	ingestion Ingestion
	commit    string
	version   string
}

// GetFields returns internal fields of a controller for testing
func (c controller) GetFields() MockDeliveryWithFields {
	return MockDeliveryWithFields{
		delivery:  c.delivery,
		ingestion: c.ingestion,
		commit:    c.commit,
		version:   c.version,
	}
}
