package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gosh/ssh-conn/db"

	"github.com/kevinburke/ssh_config"
	"github.com/spf13/cobra"
)

var tag string

// 彩色输出标签
func colorTag(tag string) string {
	colors := []string{"\033[36m", "\033[32m", "\033[35m", "\033[33m", "\033[31m"}
	return colors[len(tag)%len(colors)] + tag + "\033[0m"
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 ssh config 主机，可按标签筛选",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := os.ExpandEnv("$HOME/.ssh/config")
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
		for _, host := range cfg.Hosts {
			if len(host.Patterns) == 0 {
				continue
			}
			name := host.Patterns[0].String()
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
			var coloredTags []string
			for _, t := range tags {
				coloredTags = append(coloredTags, colorTag(t))
			}
			if len(coloredTags) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]\n", name, strings.Join(coloredTags, ", "))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", "按标签前缀筛选主机")
}
