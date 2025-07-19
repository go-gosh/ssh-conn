package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTagAndListCommands(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(home+"/.ssh", 0700)
	os.MkdirAll(home+"/.ssh-conn", 0700)
	host := "cmdhost"
	configPath := home + "/.ssh/config"
	content := "Host cmdhost\n    HostName 1.1.1.1\n    User user\n    Port 22"
	content = strings.ReplaceAll(content, "\\n", "\n")
	os.WriteFile(configPath, []byte(content), 0600)
	buf := new(bytes.Buffer)

	// tag add
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "add", "foo", host})
	rootCmd.Execute()
	addOut := buf.String()
	if !strings.Contains(addOut, "添加标签成功") {
		t.Fatalf("tag add output error: %s", addOut)
	}
	buf.Reset()

	// tag ls
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "ls"})
	rootCmd.Execute()
	lsOut := buf.String()
	if !strings.Contains(lsOut, "foo") || !strings.Contains(lsOut, host) {
		t.Fatalf("ls should show all tags and hosts, got %s", lsOut)
	}
	buf.Reset()

	// list
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list"})
	rootCmd.Execute()
	listOut := buf.String()
	if !strings.Contains(listOut, host) || !strings.Contains(listOut, "foo") {
		t.Fatalf("list output error: %s", listOut)
	}
	buf.Reset()

	// list --tag foo
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list", "--tag", "foo"})
	rootCmd.Execute()
	tagListOut := buf.String()
	if !strings.Contains(tagListOut, host) {
		t.Fatalf("list tag filter error: %s", tagListOut)
	}
	buf.Reset()

	// tag rm
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "rm", "foo", host})
	rootCmd.Execute()
	rmOut := buf.String()
	if !strings.Contains(rmOut, "移除标签成功") {
		t.Fatalf("tag rm output error: %s", rmOut)
	}
	buf.Reset()

	// tag ls again
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "ls"})
	rootCmd.Execute()
	lsOut2 := buf.String()
	if strings.Contains(lsOut2, host) {
		t.Fatalf("rm should remove host from tag, got %s", lsOut2)
	}
}
