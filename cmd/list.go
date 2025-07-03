package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-gosh/ssh-conn/db"

	"github.com/kevinburke/ssh_config"
	"github.com/spf13/cobra"
)

var (
	tag    string
	sortBy string
	order  string
)

// 彩色输出标签
func colorTag(tag string) string {
	colors := []string{"\033[36m", "\033[32m", "\033[35m", "\033[33m", "\033[31m"}
	return colors[len(tag)%len(colors)] + tag + "\033[0m"
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
		type hostInfo struct {
			Name     string
			HostName string
			Tags     []string
		}
		var hosts []hostInfo
		for _, host := range cfg.Hosts {
			if len(host.Patterns) == 0 {
				continue
			}
			name := host.Patterns[0].String()
			hostname := ""
			for _, node := range host.Nodes {
				if kv, ok := node.(*ssh_config.KV); ok && strings.EqualFold(kv.Key, "HostName") {
					hostname = kv.Value
					break
				}
			}
			tags, _ := db.GetTags(name)
			// 标签前缀匹配
			if tag != "" {
				matched := false
				for _, t := range tags {
					if strings.HasPrefix(t, tag) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
			hosts = append(hosts, hostInfo{Name: name, HostName: hostname, Tags: tags})
		}
		// 排序
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", "按标签前缀筛选主机")
	listCmd.Flags().StringVar(&sortBy, "sort", "host", "排序字段，可选 host/hostname，默认host")
	listCmd.Flags().StringVar(&order, "order", "asc", "排序方式，可选 asc/desc，默认asc")
}
