package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundingBoxFromWKT(t *testing.T) {
	tests := []struct {
		name     string
		wkt      string
		expected BBox
	}{
		{
			name: "Point",
			wkt:  "POINT(527459.24 176380.37)",
			expected: BBox{
				MinX: 527459.24, MaxX: 527459.24,
				MinY: 176380.37, MaxY: 176380.37,
			},
		},
		{
			name: "LineString",
			wkt:  "LINESTRING(526977.674310138 181798.936219104,528476.982595162 179982.126895356,528476.982595162 179990.946626588)",
			expected: BBox{
				MinX: 526977.674310138, MaxX: 528476.982595162,
				MinY: 179982.126895356, MaxY: 181798.936219104,
			},
		},
		{
			name: "Polygon",
			wkt:  "POLYGON((0 0, 10 0, 10 5, 0 5, 0 0))",
			expected: BBox{
				MinX: 0, MaxX: 10,
				MinY: 0, MaxY: 5,
			},
		},
		{
			name: "Polygon Z",
			wkt:  "POLYGON Z ((519486.91227482 257914.829213376 0,519486.371896552 257906.61227976 0,519503.134927786 257904.454717125 0,519504.067150431 257912.782373436 0,519486.91227482 257914.829213376 0))",
			expected: BBox{
				MinX: 519486.371896552, MaxX: 519504.067150431,
				MinY: 257904.454717125, MaxY: 257914.829213376,
			},
		},
	}

	const tol = 1e-9 // tolerance for float comparisons

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbox, err := BoundingBoxFromWKT(tt.wkt)
			require.NoError(t, err)
			assert.True(t, bbox.Equals(tt.expected, tol), "got %+v, want %+v", bbox, tt.expected)
		})
	}
}

func TestBoundingBoxFromWKT_Error(t *testing.T) {
	_, err := BoundingBoxFromWKT("INVALID WKT")
	assert.Error(t, err)
}

func TestBoundingBoxFromCSV(t *testing.T) {
	tests := []struct {
		name     string
		csv      string
		expected BBox
	}{
		{
			name: "Standard",
			csv:  "0, 0, 10, 5",
			expected: BBox{
				MinX: 0, MaxX: 10,
				MinY: 0, MaxY: 5,
			},
		},
		{
			name: "Reversed",
			csv:  "10, 5, 0, 0",
			expected: BBox{
				MinX: 0, MaxX: 10,
				MinY: 0, MaxY: 5,
			},
		},
	}

	const tol = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbox, err := BoundingBoxFromCSV(tt.csv)
			require.NoError(t, err)
			assert.True(t, bbox.Equals(tt.expected, tol), "got %+v, want %+v", bbox, tt.expected)
		})
	}
}

func TestBoundingBoxFromCSV_Error(t *testing.T) {
	tests := []struct {
		name string
		csv  string
	}{
		{
			name: "Too few values",
			csv:  "0,0,0",
		},
		{
			name: "Too many values",
			csv:  "0,0,0,0,0",
		},
		{
			name: "Invalid float",
			csv:  "0,0,0,abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BoundingBoxFromCSV(tt.csv)
			assert.Error(t, err)
		})
	}
}
