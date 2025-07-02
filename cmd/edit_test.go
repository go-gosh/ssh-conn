package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestEditCommand(t *testing.T) {
	home := os.TempDir()
	os.Setenv("HOME", home)
	os.MkdirAll(home+"/.ssh", 0700)
	configPath := home + "/.ssh/config"
	content := `Host edittest
    HostName 1.1.1.1
    User old
    Port 22`
	ioutil.WriteFile(configPath, []byte(content), 0600)

	editHost = "edittest"
	editHostname = "2.2.2.2"
	editUser = "newuser"
	editPort = "2200"
	editCmd.Run(editCmd, nil)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	out := string(data)
	if !strings.Contains(out, "HostName 2.2.2.2") || !strings.Contains(out, "User newuser") || !strings.Contains(out, "Port 2200") {
		t.Fatalf("edit failed: %s", out)
	}
}
