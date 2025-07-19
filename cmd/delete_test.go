package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDeleteCommand(t *testing.T) {
	home := os.TempDir()
	os.Setenv("HOME", home)
	os.MkdirAll(home+"/.ssh", 0700)
	configPath := home + "/.ssh/config"
	content := "Host deltest\n    HostName 1.1.1.1\n    User user\n    Port 22"
	content = strings.ReplaceAll(content, "\\n", "\n")
	ioutil.WriteFile(configPath, []byte(content), 0600)

	delHost = "deltest"
	deleteCmd.Run(deleteCmd, nil)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	out := string(data)
	if strings.Contains(out, "Host deltest") {
		t.Fatalf("delete failed: %s", out)
	}
}
