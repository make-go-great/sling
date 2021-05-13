package sling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New(nil)
	assert.NotNil(t, s.header)
	assert.NotNil(t, s.queries)
}

func TestSlingClone(t *testing.T) {
	s := New(nil)
	clonedS, err := s.Clone()
	assert.NoError(t, err)
	assert.Equal(t, s.method, clonedS.method)
	assert.Equal(t, s.reqURL, clonedS.reqURL)
	assert.Equal(t, s.header, clonedS.header)
	assert.Equal(t, s.queries, clonedS.queries)
}
