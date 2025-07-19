package cmd

import (
	"sort"

	"github.com/go-gosh/ssh-conn/internal/tags"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "主机标签管理",
}

var tagAddCmd = &cobra.Command{
	Use:   "add <tag> <host1 host2 ...>",
	Short: "为主机添加标签",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		hosts := args[1:]
		err := tags.AddTagsBatch(hosts, []string{tag})
		if err != nil {
			cmd.Println("添加标签失败:", err)
			return
		}
		cmd.Println("添加标签成功")
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "rm <tag> <host1 host2 ...>",
	Short: "移除主机标签",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		hosts := args[1:]
		err := tags.RemoveTagsBatch(hosts, []string{tag})
		if err != nil {
			cmd.Println("移除标签失败:", err)
			return
		}
		cmd.Println("移除标签成功")
	},
}

var tagListCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有标签及主机",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tags.EnsureTagsFile()
		if err != nil {
			cmd.Println("无法初始化标签文件:", err)
			return
		}
		m, err := tags.LoadAll()
		if err != nil {
			cmd.Println("读取标签失败:", err)
			return
		}
		var tagNames []string
		for tag := range m {
			tagNames = append(tagNames, tag)
		}
		sort.Strings(tagNames)
		for _, tag := range tagNames {
			cmd.Printf("%s: %v\n", tag, m[tag])
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagListCmd)
}

// 批量添加标签
func AddTagsBatch(hosts, tagsList []string) error {
	for _, host := range hosts {
		err := tags.AddTags(host, tagsList)
		if err != nil {
			return err
		}
	}
	return nil
}

// 批量移除标签
func RemoveTagsBatch(hosts, tagsList []string) error {
	for _, host := range hosts {
		err := tags.RemoveTags(host, tagsList)
		if err != nil {
			return err
		}
	}
	return nil
}
