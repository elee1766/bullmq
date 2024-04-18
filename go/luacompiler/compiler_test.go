package luacompiler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultScripts(t *testing.T) {
	scripts := LoadDefaultScripts()
	s, err := scripts.Load("addDelayedJob-6")
	require.NoError(t, err)
	require.False(t, bytes.Contains(s, []byte("--- @include")))
}
