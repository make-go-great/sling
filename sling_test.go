package sling

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New(nil)
	assert.NotNil(t, s.header)
	assert.NotNil(t, s.queries)
}

func TestSlingClone(t *testing.T) {
	tests := []struct {
		name string
		s    *Sling
	}{
		{
			name: "method",
			s: &Sling{
				method:  "method",
				header:  make(http.Header),
				queries: []interface{}{},
			},
		},
		{
			name: "header",
			s: &Sling{
				header: http.Header{
					"a": {"1", "2"},
				},
				queries: []interface{}{},
			},
		},
		{
			name: "queries",
			s: &Sling{
				header: make(http.Header),
				queries: []interface{}{
					"a", 1, true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clonedS, err := tc.s.Clone()
			assert.NoError(t, err)
			assert.Equal(t, tc.s.method, clonedS.method)
			assert.Equal(t, tc.s.reqURL, clonedS.reqURL)
			assert.Equal(t, tc.s.header, clonedS.header)
			assert.Equal(t, tc.s.queries, clonedS.queries)
			assert.Equal(t, tc.s.bodyProvider, clonedS.bodyProvider)
			assert.Equal(t, tc.s.rspDecoder, clonedS.rspDecoder)
		})
	}
}
