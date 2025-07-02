package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestTagCommands(t *testing.T) {
	home := os.TempDir() + "/tagtest"
	os.Setenv("HOME", home)
	os.RemoveAll(home)
	os.MkdirAll(home+"/.ssh", 0700)
	defer os.RemoveAll(home)
	host := "cmdhost"
	configPath := home + "/.ssh/config"
	content := `Host cmdhost
    HostName 1.1.1.1
    User user
    Port 22`
	ioutil.WriteFile(configPath, []byte(content), 0600)
	buf := new(bytes.Buffer)
	// add
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "add", host, "--tags", "foo,bar/baz"})
	rootCmd.Execute()
	// list
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "list", host})
	rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "foo") || !strings.Contains(out, "bar/baz") {
		t.Fatalf("list should show all tags, got %s", out)
	}
	// edit
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "edit", host, "--tags", "dev/ops"})
	rootCmd.Execute()
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "list", host})
	rootCmd.Execute()
	out = buf.String()
	if !strings.Contains(out, "dev/ops") {
		t.Fatalf("edit should overwrite tags, got %s", out)
	}
	// remove
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "remove", host, "--tags", "dev/ops"})
	rootCmd.Execute()
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "list", host})
	rootCmd.Execute()
	out = buf.String()
	if strings.Contains(out, "dev/ops") {
		t.Fatalf("remove should delete tag, got %s", out)
	}
	// 再添加 foo，保证 all 能查到标签
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "add", host, "--tags", "foo"})
	rootCmd.Execute()
	// all
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tag", "all"})
	rootCmd.Execute()
	out = buf.String()
	if !strings.Contains(out, "foo") {
		t.Fatalf("all should list all tags, got %s", out)
	}
}
