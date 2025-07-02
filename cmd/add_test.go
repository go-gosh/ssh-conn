package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAddCommand(t *testing.T) {
	home := os.TempDir()
	os.Setenv("HOME", home)
	configPath := home + "/.ssh/config"
	os.MkdirAll(home+"/.ssh", 0700)
	os.Remove(configPath)

	addHost = "addhost"
	addHostname = "1.2.3.4"
	addUser = "tester"
	addPort = "2222"
	addTag = "testtag"
	addCmd.Run(addCmd, nil)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "Host addhost") || !strings.Contains(content, "HostName 1.2.3.4") || !strings.Contains(content, "User tester") || !strings.Contains(content, "Port 2222") {
		t.Fatalf("config not written correctly: %s", content)
	}
}
