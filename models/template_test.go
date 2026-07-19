package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// exerciseByName pulls a single exercise row out of a template by name, so the
// assertions below do not depend on sort order.
func exerciseByName(t *testing.T, exercises []*TemplateExercise, name string) *TemplateExercise {
	t.Helper()
	for _, ex := range exercises {
		if ex.Name == name {
			return ex
		}
	}
	t.Fatalf("exercise %q not found in template", name)
	return nil
}

// UpdateTemplate deletes and re-inserts every exercise row, so any field it
// forgets to copy is silently cleared on every edit. Both flags are asserted:
// a test that only checked the time-based exercise would also pass if the
// update forced every row to true.
func TestUpdateTemplate_PreservesExerciseType(t *testing.T) {
	exercises := []TemplateExerciseInput{
		{Name: "test_plank", IsTimeBased: true, Block: "abs", SortOrder: 0},
		{Name: "test_bench_press", Block: "main", SortOrder: 1},
	}

	tmpl, err := CreateTemplate("test_upper_a", "chest", exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	_, err = UpdateTemplate(tmpl.ID, "test_upper_a_renamed", "chest", exercises)
	require.NoError(t, err)

	_, saved, err := GetTemplateByID(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 2)

	plank := exerciseByName(t, saved, "test_plank")
	assert.True(t, plank.IsTimeBased, "time-based flag must survive an update")

	bench := exerciseByName(t, saved, "test_bench_press")
	assert.False(t, bench.IsTimeBased, "a weighted exercise must not become time-based")
}
