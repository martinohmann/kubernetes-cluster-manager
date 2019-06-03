package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	assert "gopkg.in/go-playground/assert.v1"
)

func helperCommandContext(ctx context.Context, s ...string) (cmd *exec.Cmd) {
	args := []string{"-test.run=TestHelperProcess", "--"}
	args = append(args, s...)

	if ctx != nil {
		cmd = exec.CommandContext(ctx, os.Args[0], args...)
	} else {
		cmd = exec.Command(os.Args[0], args...)
	}

	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	return cmd
}

func helperCommand(s ...string) *exec.Cmd {
	return helperCommandContext(nil, s...)
}

func TestRunWithContext(t *testing.T) {
	cmd := helperCommand("echo", "yay")

	out, err := RunWithContext(context.Background(), cmd)

	require.NoError(t, err)

	assert.Equal(t, "yay\n", out)
}

func TestRunSilently(t *testing.T) {
	cmd := helperCommand("echo", "foo")

	out, err := RunSilently(cmd)

	require.NoError(t, err)

	assert.Equal(t, "foo\n", out)
}

func TestRunError(t *testing.T) {
	cmd := helperCommand("nonexistent-command")

	out, err := Run(cmd)

	require.Error(t, err)

	assert.Equal(t, `unknown command "nonexistent-command"`+"\n", out)
}

func TestRunSilentlyError(t *testing.T) {
	cmd := helperCommand("nonexistent-command")

	out, err := RunSilently(cmd)

	require.Error(t, err)

	assert.Equal(t, `unknown command "nonexistent-command"`+"\n", out)
}

func TestRunSilentlyWithContextCancelAfter(t *testing.T) {
	cmd := helperCommand("echo", "bar")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, err := RunSilentlyWithContext(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, "bar\n", out)
}

func TestCancelRunSilentlyWithContext(t *testing.T) {
	ctx := context.Background()

	testRunSilentlyWithContextCancel(t, ctx, "interrupt", "SIGINT received\n")
}

func TestCancelRunSilentlyWithContextSignal(t *testing.T) {
	ctx := context.WithValue(context.Background(), CancelSignal, syscall.SIGTERM)

	testRunSilentlyWithContextCancel(t, ctx, "terminated", "SIGTERM received\n")
}

type nopExecutor struct{}

func (nopExecutor) Run(*exec.Cmd) (string, error)                                     { return "", nil }
func (nopExecutor) RunWithContext(context.Context, *exec.Cmd) (string, error)         { return "", nil }
func (nopExecutor) RunSilently(*exec.Cmd) (string, error)                             { return "", nil }
func (nopExecutor) RunSilentlyWithContext(context.Context, *exec.Cmd) (string, error) { return "", nil }

func TestRestoreExecutor(t *testing.T) {
	initial := DefaultExecutor
	custom := &nopExecutor{}

	restore := SetExecutorWithRestore(custom)

	require.Equal(t, custom, DefaultExecutor)

	restore()

	require.Equal(t, initial, DefaultExecutor)
}

func testRunSilentlyWithContextCancel(t *testing.T, ctx context.Context, c string, expected string) {
	cmd := helperCommand(c)

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)

	var out string
	var err error

	go func(wg *sync.WaitGroup) {
		out, err = RunSilentlyWithContext(ctx, cmd)

		wg.Done()
	}(&wg)

	time.Sleep(100 * time.Millisecond)

	cancel()

	wg.Wait()

	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

// This is not an actual test. It's a helper process that gets called by the
// command executor tests. This is the same approach as in the tests for the
// os/exec package: https://github.com/golang/go/blob/master/src/os/exec/exec_test.go
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args

	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]

	switch cmd {
	case "echo":
		iargs := []interface{}{}
		for _, s := range args {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)
	case "interrupt":
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)

		select {
		case <-signalChan:
			fmt.Println("SIGINT received")
		case <-time.After(500 * time.Millisecond):
			fmt.Println("timeout")
		}
	case "terminated":
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM)

		select {
		case <-signalChan:
			fmt.Println("SIGTERM received")
		case <-time.After(500 * time.Millisecond):
			fmt.Println("timeout")
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", cmd)
		os.Exit(2)
	}
}
