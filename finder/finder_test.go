package finder_test

import (
	"testing"

	"github.com/skyaxl/golactivity/finder"
	"github.com/skyaxl/golactivity/reader"

	"github.com/stretchr/testify/require"
)

func TestReadTokens(t *testing.T) {
	r := reader.New("./../example", ".*")
	pgks, fset, _ := r.Read()
	funcs := finder.FindAllowedFuncs(pgks, fset)
	require.Len(t, funcs, 1)
	require.Equal(t, "Execute", funcs[0].Name.Name)
}
