package regression

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepositories(t *testing.T) {
	require := require.New(t)

	tmpDir, err := createTempDir()
	require.NoError(err)
	defer os.RemoveAll(tmpDir)

	config := NewConfig()
	require.NotNil(config)

	config.RepositoriesCache = tmpDir
	config.Complexity = 0

	repos := NewDefaultRepositories(config)
	require.NotNil(repos)
	require.Equal(config, repos.config)
	require.Equal(tmpDir, repos.Path())

	err = repos.Download()
	require.NoError(err)

	r, err := ioutil.ReadDir(tmpDir)
	require.NoError(err)

	linkDir, err := repos.LinksDir()
	require.NoError(err)
	defer os.RemoveAll(linkDir)

	links, err := ioutil.ReadDir(linkDir)
	require.NoError(err)
	require.Len(links, len(r))

	for i, link := range links {
		require.Equal(r[i].Name(), link.Name())
	}
}
