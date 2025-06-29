package orderedmap_test

import (
	"orderedmap/internal/orderedmap"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	t.Run("Add and Get", func(t *testing.T) {
		om := orderedmap.NewOrderedMap()
		om.Add("name", "John")

		val, exists := om.Get("name")
		assert.True(t, exists)
		assert.Equal(t, "John", val)
	})

	t.Run("Update existing key", func(t *testing.T) {
		om := orderedmap.NewOrderedMap()
		om.Add("name", "John")
		om.Add("name", "Jane")

		val, exists := om.Get("name")
		assert.True(t, exists)
		assert.Equal(t, "Jane", val)
	})

	t.Run("Delete key", func(t *testing.T) {
		om := orderedmap.NewOrderedMap()
		om.Add("name", "John")
		om.Delete("name")

		_, exists := om.Get("name")
		assert.False(t, exists)
	})

}
