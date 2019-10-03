package regression

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getGopath() (string, error) {
	gopath := os.Getenv("GOPATH")
	split := strings.Split(gopath, ":")

	if len(split) > 0 {
		return split[0], nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, "go"), nil
}

func TestServer(t *testing.T) {
	require := require.New(t)

	gopath, err := getGopath()
	require.NoError(err)

	config := GitServerConfig{
		RepositoriesCache: fmt.Sprintf("%s/src/github.com/src-d", gopath),
		GitServerPort:     9418,
	}

	server := NewGitServer(config)

	err = server.Start()
	require.NoError(err)
	require.True(server.Alive())

	err = server.Stop()
	require.NoError(err)
	require.False(server.Alive())

	rusage := server.Rusage()
	require.NotNil(rusage)
	require.True(rusage.Stime.Nano()+rusage.Utime.Nano() > 0)
	require.True(rusage.Maxrss > 0)
}
