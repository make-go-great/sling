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
