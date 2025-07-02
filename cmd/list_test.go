package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/go-gosh/ssh-conn/db"
)

func TestListCommand(t *testing.T) {
	home := os.TempDir() + "/listtest"
	os.Setenv("HOME", home)
	os.RemoveAll(home)
	os.MkdirAll(home+"/.ssh", 0700)
	defer os.RemoveAll(home)
	configPath := home + "/.ssh/config"
	content := `Host lister
    HostName 1.1.1.1
    User user
    Port 22`
	ioutil.WriteFile(configPath, []byte(content), 0600)
	// 添加标签
	tags, _ := db.GetTags("lister")
	if len(tags) == 0 {
		db.AddTags("lister", []string{"dev/web", "ops"})
	}
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list"})
	rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "lister") || !strings.Contains(out, "dev/web") || !strings.Contains(out, "ops") {
		t.Fatalf("list output error: %s", out)
	}
	// 前缀筛选
	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list", "--tag", "dev/"})
	rootCmd.Execute()
	out = buf.String()
	if !strings.Contains(out, "dev/web") {
		t.Fatalf("list tag prefix filter error: %s", out)
	}
}
