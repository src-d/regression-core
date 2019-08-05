package regression

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinary(t *testing.T) {
	require := require.New(t)

	dir, err := CreateTempDir()
	require.NoError(err)
	defer os.RemoveAll(dir)

	tool, releases := setupBinary(t)

	config := NewConfig()
	config.BinaryCache = dir

	cases := []struct {
		version  string
		binSize  int64
		yamlSize int64
		ok       bool
	}{
		{
			version:  "v0.23.1",
			binSize:  42596563,
			yamlSize: 5007,
			ok:       true,
		},
		{
			version:  "remote:v0.16.0",
			binSize:  0,
			yamlSize: 5072,
			ok:       true,
		},
	}

	for _, c := range cases {
		err = os.RemoveAll(dir)
		require.NoError(err)

		binary := NewBinary(config, tool, c.version, releases)
		err = binary.Download()

		if !c.ok {
			require.Error(err)
			continue
		}

		require.NoError(err)

		if c.binSize > 0 {
			s, err := os.Stat(binary.Path)
			require.NoError(err)

			require.Equal(c.binSize, s.Size())
		}

		if c.yamlSize > 0 {
			s, err := os.Stat(binary.ExtraFile("regression.yml"))
			require.NoError(err)

			require.Equal(c.yamlSize, s.Size())
		}
	}
}

func setupBinary(t *testing.T) (Tool, *Releases) {
	t.Helper()

	token := os.Getenv("REG_TOKEN")
	if token == "" {
		t.Skip("REG_TOKEN not provided")
	}

	tool := Tool{
		Name:        "gitbase",
		GitURL:      "https://github.com/src-d/gitbase",
		ProjectPath: "github.com/src-d/gitbase",
		BuildSteps: []BuildStep{
			{
				Dir:     "",
				Command: "make",
				Args:    []string{"dependencies", "packages"},
			},
		},
		ExtraFiles: []string{
			"_testdata/regression.yml",
		},
	}

	releases := NewReleases("src-d", "gitbase", token)

	return tool, releases
}
