package regression

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	require := require.New(t)

	version := "remote:v0.12.1"

	step := BuildStep{
		Dir:     "",
		Command: "make",
		Args:    []string{"packages"},
	}

	tool := Tool{
		Name:        "borges",
		GitURL:      "https://github.com/src-d/borges",
		ProjectPath: "github.com/src-d/borges",
		BuildSteps: []BuildStep{
			step,
		},
	}

	build, err := NewBuild(NewConfig(), tool, version)
	require.NoError(err)

	_, err = build.download()
	require.NoError(err)

	err = build.buildStep(step)
	require.NoError(err)
}
