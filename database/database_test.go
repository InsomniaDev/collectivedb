package database

// import (
// 	"fmt"
// 	"testing"

// 	assert "github.com/stretchr/testify/assert"
// )

// func TestDatabase_Init(t *testing.T) {
// 	// Test the initialization
// 	resp := Init(1)
// 	assert.Equal(t, 1, resp.TotalDepth)
// }

// func TestDatabase_FNV32a(t *testing.T) {
// 	// Test converting from string through a hash to an int
// 	assert.Equal(t, "2949673445", fmt.Sprint(FNV32a("test")))
// }

// func TestDatabase_HashNum(t *testing.T) {
// 	// Test returning a num with length of three
// 	respInt, respStr := hashNum("test", 3, "")
// 	assert.Equal(t, 294, respInt)
// 	assert.Equal(t, "2949673445", respStr)

// 	// Test with the depth at zero
// 	respInt, respStr = hashNum("test", 0, "2949673445")
// 	assert.Equal(t, 0, respInt)
// 	assert.Equal(t, "2949673445", respStr)

// 	// Test with a hast string provided
// 	respInt, respStr = hashNum("test", 3, "2949673445")
// 	assert.Equal(t, 294, respInt)
// 	assert.Equal(t, "2949673445", respStr)
// }
