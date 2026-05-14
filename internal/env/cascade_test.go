package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCascadeBasicOrder(t *testing.T) {
	layers := map[string]map[string]string{
		"base":    {"A": "1", "B": "2"},
		"staging": {"B": "20", "C": "3"},
	}
	opts := CascadeOptions{Layers: []string{"base", "staging"}, Overwrite: true}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)
	assert.Equal(t, "1", res.Vars["A"])
	assert.Equal(t, "20", res.Vars["B"])
	assert.Equal(t, "3", res.Vars["C"])
	assert.Equal(t, "base", res.Origin["A"])
	assert.Equal(t, "staging", res.Origin["B"])
}

func TestCascadeNoOverwrite(t *testing.T) {
	layers := map[string]map[string]string{
		"base":    {"A": "base-val"},
		"overlay": {"A": "overlay-val"},
	}
	opts := CascadeOptions{Layers: []string{"base", "overlay"}, Overwrite: false}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)
	assert.Equal(t, "base-val", res.Vars["A"])
	assert.Equal(t, "base", res.Origin["A"])
}

func TestCascadeMissingLayer(t *testing.T) {
	layers := map[string]map[string]string{
		"base": {"A": "1"},
	}
	opts := CascadeOptions{Layers: []string{"base", "missing"}, Overwrite: true}
	_, err := Cascade(layers, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
}

func TestCascadeSummaryLines(t *testing.T) {
	layers := map[string]map[string]string{
		"base": {"X": "1", "Y": "2"},
	}
	opts := CascadeOptions{Layers: []string{"base"}, Overwrite: false}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)
	lines := res.SummaryLines()
	assert.Len(t, lines, 2)
	assert.Contains(t, lines[0], "base")
}

func TestCascadeEmptyLayers(t *testing.T) {
	opts := CascadeOptions{Layers: []string{}, Overwrite: true}
	res, err := Cascade(map[string]map[string]string{}, opts)
	require.NoError(t, err)
	assert.Empty(t, res.Vars)
}

func TestCascadeOriginTracking(t *testing.T) {
	layers := map[string]map[string]string{
		"dev":  {"DB_URL": "dev-db"},
		"prod": {"DB_URL": "prod-db", "SECRET": "s3cr3t"},
	}
	opts := CascadeOptions{Layers: []string{"dev", "prod"}, Overwrite: true}
	res, err := Cascade(layers, opts)
	require.NoError(t, err)
	assert.Equal(t, "prod", res.Origin["DB_URL"])
	assert.Equal(t, "prod", res.Origin["SECRET"])
}
