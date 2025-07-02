package cmd

import (
	"os/exec"
	"testing"
)

func TestConnectCommand(t *testing.T) {
	called := false
	sshExec = func(name string, arg ...string) *exec.Cmd {
		called = true
		return exec.Command("echo", "mock ssh")
	}
	defer func() { sshExec = exec.Command }()

	rootCmd.SetArgs([]string{"connect", "testhost"})
	rootCmd.Execute()

	if !called {
		t.Fatal("sshExec should be called")
	}
}
