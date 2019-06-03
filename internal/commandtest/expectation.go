package commandtest

import "fmt"

// ExpectedCommand is used to configure expectations a the MockExecutor.
type ExpectedCommand struct {
	fulfilled bool
	execute   bool
	err       error
	out       string
	command   string
}

// WillReturnOutput will configure the ExpectedCommand to return out.
func (e *ExpectedCommand) WillReturn(out string) *ExpectedCommand {
	e.out = out
	return e
}

// WillReturnError will configure the ExpectedCommand to return error err.
func (e *ExpectedCommand) WillReturnError(err error) *ExpectedCommand {
	e.err = err
	return e
}

// WillExecute marks the expected command to be run with a command.Executor
// wrapped by the MockExecutor.
func (e *ExpectedCommand) WillExecute() *ExpectedCommand {
	e.execute = true
	return e
}

// String implements fmt.Stringer.
func (e *ExpectedCommand) String() string {
	msg := fmt.Sprintf("command %q", e.command)
	if e.execute {
		msg += ", which should execute"
	} else if e.out != "" || e.err != nil {
		msg += ", which should:"
		if e.out != "" {
			msg += fmt.Sprintf("\n- return output: %s", e.out)
		}

		if e.err != nil {
			msg += fmt.Sprintf("\n- return error: %s", e.err)
		}
	}

	return msg
}
