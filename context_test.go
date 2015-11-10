package xid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestFromContext(t *testing.T) {
	id, ok := FromContext(nil)
	assert.False(t, ok)
	assert.Equal(t, ID{}, id)
	id, ok = FromContext(context.Background())
	assert.False(t, ok)
	assert.Equal(t, ID{}, id)
	id2 := New()
	ctx := NewContext(context.Background(), id2)
	id, ok = FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, id2, id)
}
