package internal

import (
	"testing"

	"github.com/map-services/street-manager-relay/models"
	"github.com/stretchr/testify/assert"
)

func TestEventMapping(t *testing.T) {
	event := &models.Event{}
	
	pointers := eventPointers(event)
	values := eventValues(event)

	// Ensure the number of columns, pointers, and values are all in sync
	assert.Equal(t, len(eventColumns), len(pointers), "Column count mismatch with pointers")
	assert.Equal(t, len(eventColumns), len(values), "Column count mismatch with values")

	// Basic check that a value can be retrieved
	foo := "bar"
	event.EventType = foo
	values = eventValues(event)
	
	// Value at index 0 should be EventType
	assert.Equal(t, foo, values[0])
	
	// Updating via pointer should update the struct
	p := pointers[0].(*string)
	*p = "baz"
	assert.Equal(t, "baz", event.EventType)
}
