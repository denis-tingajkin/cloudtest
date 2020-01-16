package runners

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/denis-tingajkin/cloudtest/pkg/model"
	"github.com/denis-tingajkin/cloudtest/pkg/shell"
	"github.com/denis-tingajkin/cloudtest/pkg/utils"
)

type shellTestRunner struct {
	test        *model.TestEntry
	envMgr      shell.EnvironmentManager
	artifactDir string
	id          string
}

func (runner *shellTestRunner) Run(timeoutCtx context.Context, env []string, writer *bufio.Writer) error {
	err := runner.runCmd(timeoutCtx, utils.ParseScript(runner.test.RunScript), env, writer)
	return err
}

func (runner *shellTestRunner) runCmd(context context.Context, script, env []string, writer *bufio.Writer) error {
	for _, cmd := range script {
		if strings.TrimSpace(cmd) == "" {
			continue
		}

		cmdEnv := append(runner.envMgr.GetProcessedEnv(), env...)
		_, _ = writer.WriteString(fmt.Sprintf(">>>>>>Running: %s:<<<<<<\n", cmd))
		_ = writer.Flush()

		logger := func(s string) {
		}
		_, err := utils.RunCommand(context, cmd, "", logger, writer, cmdEnv, map[string]string{"artifacts-dir": runner.artifactDir}, false)
		if err != nil {
			_, _ = writer.WriteString(fmt.Sprintf("error running command: %v\n", err))
			_ = writer.Flush()
			return err
		}
	}
	return nil
}

func (runner *shellTestRunner) GetCmdLine() string {
	return runner.test.RunScript
}

// NewShellTestRunner - creates a new shell script test runner.
func NewShellTestRunner(ids string, test *model.TestEntry) TestRunner {
	envMgr := shell.NewEnvironmentManager()
	_ = envMgr.ProcessEnvironment(ids, "shellrun", os.TempDir(), test.ExecutionConfig.Env, map[string]string{})
	artifactDir := ""
	if len(test.ArtifactDirectories) > 0 {
		artifactDir = test.ArtifactDirectories[len(test.ArtifactDirectories)-1]
	}
	return &shellTestRunner{
		id:          ids,
		test:        test,
		envMgr:      envMgr,
		artifactDir: artifactDir,
	}
}
