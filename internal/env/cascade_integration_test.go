package env

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCascadeRoundtripWithParser(t *testing.T) {
	baseRaw := "APP_ENV=base\nDB_HOST=localhost\n"
	prodRaw := "APP_ENV=production\nDB_HOST=prod.db.internal\nSECRET_KEY=abc123\n"

	baseVars, err := Parse(strings.NewReader(baseRaw))
	require.NoError(t, err)
	prodVars, err := Parse(strings.NewReader(prodRaw))
	require.NoError(t, err)

	layers := map[string]map[string]string{
		"base": baseVars,
		"prod": prodVars,
	}
	opts := CascadeOptions{Layers: []string{"base", "prod"}, Overwrite: true}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)

	assert.Equal(t, "production", res.Vars["APP_ENV"])
	assert.Equal(t, "prod.db.internal", res.Vars["DB_HOST"])
	assert.Equal(t, "abc123", res.Vars["SECRET_KEY"])
	assert.Equal(t, "prod", res.Origin["APP_ENV"])
	assert.Equal(t, "base", res.Origin["DB_HOST"], "base origin should be overwritten")
}

func TestCascadeSerializeResult(t *testing.T) {
	layers := map[string]map[string]string{
		"base":    {"FOO": "bar", "BAZ": "qux"},
		"overlay": {"FOO": "overridden"},
	}
	opts := CascadeOptions{Layers: []string{"base", "overlay"}, Overwrite: true}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)

	var sb strings.Builder
	err = Serialize(res.Vars, &sb)
	require.NoError(t, err)

	output := sb.String()
	assert.Contains(t, output, "FOO=overridden")
	assert.Contains(t, output, "BAZ=qux")
}
