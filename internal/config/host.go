package config

import (
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"
)

type Host struct {
	Name         string
	HostName     string
	User         string
	Port         string
	IdentityFile string
	ProxyJump    string
}

func getHomeDir() string {
	h := os.Getenv("HOME")
	if h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return home
}

// ParseSSHConfig 解析 ~/.ssh/config，返回所有 Host 条目
func ParseSSHConfig() ([]Host, error) {
	home := getHomeDir()
	configPath := filepath.Join(home, ".ssh", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cfg, err := ssh_config.Decode(file)
	if err != nil {
		return nil, err
	}
	var hosts []Host
	for _, h := range cfg.Hosts {
		if len(h.Patterns) == 0 {
			continue
		}
		name := h.Patterns[0].String()
		if name == "" {
			continue
		}
		host := Host{Name: name}
		host.HostName, _ = cfg.Get(name, "HostName")
		host.User, _ = cfg.Get(name, "User")
		host.Port, _ = cfg.Get(name, "Port")
		host.IdentityFile, _ = cfg.Get(name, "IdentityFile")
		host.ProxyJump, _ = cfg.Get(name, "ProxyJump")
		hosts = append(hosts, host)
	}
	return hosts, nil
}
