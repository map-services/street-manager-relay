package internal

import (
	"os"
	"testing"
	"time"

	"github.com/map-services/street-manager-relay/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestDbRepository_Upsert(t *testing.T) {
	dbPath := "test_upsert.db"
	defer func() { _ = os.Remove(dbPath) }()

	repo, err := NewDbRepository(dbPath)
	require.NoError(t, err)
	defer func(){ _ = repo.Close() }()

	street1 := "High Street"
	coords1 := "POINT(0.1278 51.5074)"
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	event := &models.Event{
		ObjectReference:     "REF123",
		EventType:           "Permit",
		StreetName:          &street1,
		ActivityCoordinates: &coords1,
		ProposedStartDate:   &now,
		ProposedEndDate:     &tomorrow,
	}

	// First upsert
	id1, err := repo.UpsertSingle(event)
	assert.NoError(t, err)
	assert.Greater(t, id1, int64(0))

	// Update event
	street2 := "Low Street"
	event.StreetName = &street2
	id2, err := repo.UpsertSingle(event)
	assert.NoError(t, err)
	assert.Equal(t, id1, id2)

	// Verify update using Search
	bbox := &models.BBox{MinX: 0.127, MaxX: 0.128, MinY: 51.507, MaxY: 51.508}
	facets := &models.Facets{}
	temporal := &models.TemporalFilters{MaxDaysAhead: 30, MaxDaysBehind: 30}

	events, err := repo.Search(bbox, facets, temporal)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "Low Street", *events[0].StreetName)
	assert.Equal(t, "REF123", events[0].ObjectReference)
}
