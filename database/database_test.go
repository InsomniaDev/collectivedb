package database

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	resp := Init(1)
	assert.Equal(t, resp.TotalDepth, 1)
}
