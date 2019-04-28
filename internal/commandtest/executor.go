package commandtest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
)

// WithMockExecutor will replace the default executor with a
// command.MockExecutor and then call f. The executor is passed to f to be able
// to make assertions and control command return values. WithMockExecutor will
// restore the previous default executor after f returns or panics. Optionally
// another executor can be passed which will be used by the
// command.MockExecutor if passthrough of a command is explicitly requested.
func WithMockExecutor(f func(*command.MockExecutor), wrapped ...command.Executor) {
	var wrappedE command.Executor
	if len(wrapped) > 0 {
		wrappedE = wrapped[0]
	}

	executor := command.NewMockExecutor(wrappedE)
	restoreExecutor := command.SetExecutorWithRestore(executor)
	defer restoreExecutor()

	f(executor)
}
