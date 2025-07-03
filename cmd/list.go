package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/go-gosh/ssh-conn/db"

	"github.com/kevinburke/ssh_config"
	"github.com/spf13/cobra"
)

var (
	tag     string
	sortBy  string
	order   string
	connect bool
)

// 彩色输出标签
func colorTag(tag string) string {
	colors := []string{"\033[36m", "\033[32m", "\033[35m", "\033[33m", "\033[31m"}
	return colors[len(tag)%len(colors)] + tag + "\033[0m"
}

type hostInfo struct {
	Name     string
	HostName string
	Tags     []string
}

func getFilteredHosts(cfg *ssh_config.Config, tagList []string) []hostInfo {
	hosts := []hostInfo{}
	for _, host := range cfg.Hosts {
		if len(host.Patterns) == 0 {
			continue
		}
		name := host.Patterns[0].String()
		hostname, _ := cfg.Get(name, "HostName")
		tags, _ := db.GetTags(name)
		if len(tagList) > 0 {
			matched := true
			for _, want := range tagList {
				found := false
				for _, t := range tags {
					if strings.HasPrefix(t, want) {
						found = true
						break
					}
				}
				if !found {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
		}
		hosts = append(hosts, hostInfo{Name: name, HostName: hostname, Tags: tags})
	}
	return hosts
}

func sortHosts(hosts []hostInfo, sortBy, order string) {
	switch sortBy {
	case "hostname":
		sort.Slice(hosts, func(i, j int) bool {
			if order == "desc" {
				return hosts[i].HostName > hosts[j].HostName
			}
			return hosts[i].HostName < hosts[j].HostName
		})
	default:
		sort.Slice(hosts, func(i, j int) bool {
			if order == "desc" {
				return hosts[i].Name > hosts[j].Name
			}
			return hosts[i].Name < hosts[j].Name
		})
	}
}

func printHosts(cmd *cobra.Command, hosts []hostInfo) {
	for _, h := range hosts {
		sort.Strings(h.Tags)
		var coloredTags []string
		for _, t := range h.Tags {
			coloredTags = append(coloredTags, colorTag(t))
		}
		if len(coloredTags) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s) [%s]\n", h.Name, h.HostName, strings.Join(coloredTags, ", "))
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s)\n", h.Name, h.HostName)
		}
	}
}

func connectHost(name string) error {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return errors.New("未找到 ssh 命令")
	}
	c := sshExec(sshPath, name)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		return errors.New("ssh 连接失败:" + err.Error())
	}
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 ssh config 主机，可按标签筛选",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "无法获取用户主目录:", err)
			return
		}
		configPath := home + "/.ssh/config"
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println("无法打开 ssh config:", err)
			return
		}
		defer file.Close()
		cfg, err := ssh_config.Decode(file)
		if err != nil {
			fmt.Println("解析 ssh config 失败:", err)
			return
		}
		tagList := []string{}
		if tag != "" {
			tagList = strings.Split(tag, ",")
			for i := range tagList {
				tagList[i] = strings.TrimSpace(tagList[i])
			}
		}
		hosts := getFilteredHosts(cfg, tagList)
		sortHosts(hosts, sortBy, order)
		printHosts(cmd, hosts)
		if connect && len(hosts) == 1 {
			err := connectHost(hosts[0].Name)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", "按标签前缀筛选主机")
	listCmd.Flags().StringVar(&sortBy, "sort", "host", "排序字段，可选 host/hostname，默认host")
	listCmd.Flags().StringVar(&order, "order", "asc", "排序方式，可选 asc/desc，默认asc")
	listCmd.Flags().BoolVarP(&connect, "connect", "c", false, "是否连接主机")
}
