package reader_test

import (
	"testing"

	"github.com/skyaxl/golactivity/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_ReadOK(t *testing.T) {
	r := reader.New("./../example", ".*")
	require.NotNil(t, r)
	packages, fset, err := r.Read()
	require.Len(t, packages, 1)
	require.Contains(t, packages, "example")
	require.NotNil(t, fset)
	require.Nil(t, err)
}

func TestReader_ReadFailed(t *testing.T) {
	r := reader.New("./../noexists", ".*")
	require.NotNil(t, r)
	packages, fset, err := r.Read()
	require.Len(t, packages, 0)
	require.NotContains(t, packages, "example")
	require.NotNil(t, fset)
	require.EqualError(t, err, "open ./../noexists: no such file or directory")
}

func TestReader_ReadNoPackages(t *testing.T) {
	r := reader.New("./../example/nopackages", ".*")
	require.NotNil(t, r)
	packages, fset, err := r.Read()
	require.Len(t, packages, 0)
	require.NotContains(t, packages, "example")
	require.NotNil(t, fset)
	require.Nil(t, err)
}
