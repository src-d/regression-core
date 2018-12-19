package regression

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleases(t *testing.T) {
	require := require.New(t)

	token := os.Getenv("REG_TOKEN")
	if token == "" {
		t.Skip("REG_TOKEN not provided")
	}

	dir, err := CreateTempDir()
	require.NoError(err)
	defer os.RemoveAll(dir)

	r := NewReleases("src-d", "borges", token)
	require.NotNil(r)
	require.Nil(r.repoReleases)

	path := filepath.Join(dir, "invalid_version")
	_, err = r.Get("invalid_version", "invalid_asset", path)
	require.Error(err)

	if !ErrVersionNotFound.Is(err) {
		t.Errorf("the error should be invalid version but is: %v", err)
		t.FailNow()
	}

	require.False(fileExist(path))

	list := r.repoReleases

	path = filepath.Join(dir, "invalid_asset")
	_, err = r.Get("v0.12.0", "invalid_asset", path)
	require.Error(err)
	require.True(ErrAssetNotFound.Is(err))
	require.False(fileExist(path))

	require.Exactly(list, r.repoReleases)

	path = filepath.Join(dir, "borges.v0.12.0")
	_, err = r.Get("v0.12.0", "borges_v0.12.0_linux_amd64.tar.gz", path)
	require.NoError(err)
	require.True(fileExist(path))
}
