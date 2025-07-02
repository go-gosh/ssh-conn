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

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "主机标签管理",
}

var tagAddCmd = &cobra.Command{
	Use:   "add [host] --tags tag1,tag2",
	Short: "为主机添加标签（可多个）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		tags := parseTags(tagsArg)
		if len(tags) == 0 {
			fmt.Println("请指定至少一个标签")
			return
		}
		err := db.AddTags(host, tags)
		if err != nil {
			fmt.Println("添加标签失败:", err)
			return
		}
		fmt.Println("添加标签成功")
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove [host] --tags tag1,tag2",
	Short: "移除主机标签（可多个）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		tags := parseTags(tagsArg)
		if len(tags) == 0 {
			fmt.Println("请指定至少一个标签")
			return
		}
		err := db.RemoveTags(host, tags)
		if err != nil {
			fmt.Println("移除标签失败:", err)
			return
		}
		fmt.Println("移除标签成功")
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list [host]",
	Short: "列出主机所有标签",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		tags, err := db.GetTags(host)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "查询标签失败:", err)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "标签:", strings.Join(tags, ", "))
	},
}

var tagEditCmd = &cobra.Command{
	Use:   "edit [host] --tags tag1,tag2",
	Short: "覆盖主机标签（可多个）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		tags := parseTags(tagsArg)
		err := db.SetTags(host, tags)
		if err != nil {
			fmt.Println("设置标签失败:", err)
			return
		}
		fmt.Println("设置标签成功")
	},
}

var tagAllCmd = &cobra.Command{
	Use:   "all",
	Short: "列出所有标签",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tagsMap := map[string]struct{}{}
		allHosts, _ := getAllHosts()
		for _, host := range allHosts {
			tags, _ := db.GetTags(host)
			for _, t := range tags {
				tagsMap[t] = struct{}{}
			}
		}
		var tags []string
		for t := range tagsMap {
			tags = append(tags, t)
		}
		sort.Strings(tags)
		fmt.Fprintln(cmd.OutOrStdout(), "所有标签:", strings.Join(tags, ", "))
	},
}

var tagsArg string

func parseTags(s string) []string {
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// 获取所有主机名
func getAllHosts() ([]string, error) {
	configPath := os.ExpandEnv("$HOME/.ssh/config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cfg, err := ssh_config.Decode(file)
	if err != nil {
		return nil, err
	}
	var hosts []string
	for _, host := range cfg.Hosts {
		if len(host.Patterns) == 0 {
			continue
		}
		hosts = append(hosts, host.Patterns[0].String())
	}
	return hosts, nil
}

func init() {
	rootCmd.AddCommand(tagCmd)

	tagCmd.AddCommand(tagAddCmd)

	tagAddCmd.Flags().StringVar(&tagsArg, "tags", "", "标签，逗号分隔")

	tagCmd.AddCommand(tagRemoveCmd)

	tagRemoveCmd.Flags().StringVar(&tagsArg, "tags", "", "标签，逗号分隔")

	tagCmd.AddCommand(tagListCmd)

	tagCmd.AddCommand(tagEditCmd)

	tagEditCmd.Flags().StringVar(&tagsArg, "tags", "", "标签，逗号分隔")

	tagCmd.AddCommand(tagAllCmd)
}
