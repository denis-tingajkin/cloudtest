package runners

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/denis-tingajkin/cloudtest/pkg/model"
	"github.com/denis-tingajkin/cloudtest/pkg/shell"
	"github.com/denis-tingajkin/cloudtest/pkg/utils"
)

type goTestRunner struct {
	test        *model.TestEntry
	cmdLine     string
	envMgr      shell.EnvironmentManager
	artifactDir string
}

func (runner *goTestRunner) Run(timeoutCtx context.Context, env []string, writer *bufio.Writer) error {
	logger := func(s string) {}
	cmdEnv := append(runner.envMgr.GetProcessedEnv(), env...)
	_, err := utils.RunCommand(timeoutCtx, runner.cmdLine, runner.test.ExecutionConfig.PackageRoot,
		logger, writer, cmdEnv, map[string]string{"artifact-dir": runner.artifactDir}, false)
	return err
}

func (runner *goTestRunner) GetCmdLine() string {
	return runner.cmdLine
}

// NewGoTestRunner - creates go test runner
func NewGoTestRunner(ids string, test *model.TestEntry, timeout time.Duration) TestRunner {
	cmdLine := fmt.Sprintf("go test . -test.timeout %v -count 1 --run \"^(%s)$\\\\z\" --tags \"%s\" --test.v",
		timeout, test.Name, test.Tags)

	envMgr := shell.NewEnvironmentManager()
	_ = envMgr.ProcessEnvironment(ids, "gotest", os.TempDir(), test.ExecutionConfig.Env, map[string]string{})
	artifactDir := ""
	if len(test.ArtifactDirectories) > 0 {
		artifactDir = test.ArtifactDirectories[len(test.ArtifactDirectories)-1]
	}
	return &goTestRunner{
		test:        test,
		cmdLine:     cmdLine,
		envMgr:      envMgr,
		artifactDir: artifactDir,
	}
}
