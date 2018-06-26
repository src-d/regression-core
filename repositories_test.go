package regression

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepositories(t *testing.T) {
	require := require.New(t)

	tmpDir, err := CreateTempDir()
	require.NoError(err)
	defer os.RemoveAll(tmpDir)

	config := NewConfig()
	require.NotNil(config)

	config.RepositoriesCache = tmpDir
	config.Complexity = 0

	repos, err := NewRepositories(config)
	require.NotNil(repos)
	require.NoError(err)
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

var repositoriesExamples = []RepoDescription{
	{
		Name:        "name",
		URL:         "url",
		Description: "description",
		Complexity:  0,
	}, {

		Name:        "go-git",
		URL:         "git://github.com/src-d/go-git",
		Description: "go-git repository",
		Complexity:  2,
	}, {
		Name: "kernel",
		URL:  "https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git",
		Description: `very
long
description
`,
		Complexity: 10,
	},
}

func TestRepositoriesYaml(t *testing.T) {
	require := require.New(t)

	config := NewConfig()
	config.RepositoriesFile = "testdata/repositories.yaml"
	r, err := NewRepositories(config)
	require.NoError(err)
	require.NotNil(r)

	require.Equal(repositoriesExamples, r.Repos)
}
