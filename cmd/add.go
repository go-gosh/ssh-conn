package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gosh/ssh-conn/db"

	"github.com/spf13/cobra"
)

var (
	addHost     string
	addHostname string
	addUser     string
	addPort     string
	addTag      string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加主机到 ssh config，并可选打标签",
	Run: func(cmd *cobra.Command, args []string) {
		if addHost == "" || addHostname == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "--host 和 --hostname 必填")
			return
		}
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "无法获取用户主目录:", err)
			return
		}
		configPath := home + "/.ssh/config"
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "无法打开 ssh config:", err)
			return
		}
		defer f.Close()
		lines := []string{"", fmt.Sprintf("Host %s", addHost), fmt.Sprintf("    HostName %s", addHostname)}
		if addUser != "" {
			lines = append(lines, fmt.Sprintf("    User %s", addUser))
		}
		if addPort != "" {
			lines = append(lines, fmt.Sprintf("    Port %s", addPort))
		}
		_, err = f.WriteString(strings.Join(lines, "\n") + "\n")
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "写入 ssh config 失败:", err)
			return
		}
		if addTag != "" {
			tags := []string{addTag}
			db.AddTags(addHost, tags)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "添加主机成功")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addHost, "host", "", "主机别名 (Host)")
	addCmd.Flags().StringVar(&addHostname, "hostname", "", "主机地址 (HostName)")
	addCmd.Flags().StringVar(&addUser, "user", "", "用户名 (User)")
	addCmd.Flags().StringVar(&addPort, "port", "", "端口 (Port)")
	addCmd.Flags().StringVar(&addTag, "tag", "", "标签")
}
