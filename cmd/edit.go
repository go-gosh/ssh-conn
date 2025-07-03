package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	editHost     string
	editHostname string
	editUser     string
	editPort     string
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑 ssh config 中的主机配置",
	Run: func(cmd *cobra.Command, args []string) {
		if editHost == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "--host 必填")
			return
		}
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "无法获取用户主目录:", err)
			return
		}
		configPath := home + "/.ssh/config"
		input, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "无法读取 ssh config:", err)
			return
		}
		lines := strings.Split(string(input), "\n")
		var out []string
		inBlock := false
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "Host ") && strings.TrimSpace(trim[5:]) == editHost {
				inBlock = true
				out = append(out, line)
				continue
			}
			if inBlock && (strings.HasPrefix(trim, "Host ") && strings.TrimSpace(trim[5:]) != editHost) {
				inBlock = false
			}
			if inBlock {
				if editHostname != "" && strings.HasPrefix(trim, "HostName ") {
					line = "    HostName " + editHostname
				}
				if editUser != "" && strings.HasPrefix(trim, "User ") {
					line = "    User " + editUser
				}
				if editPort != "" && strings.HasPrefix(trim, "Port ") {
					line = "    Port " + editPort
				}
			}
			out = append(out, line)
		}
		os.WriteFile(configPath, []byte(strings.Join(out, "\n")), 0600)
		fmt.Fprintln(cmd.OutOrStdout(), "编辑完成")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVar(&editHost, "host", "", "主机别名 (Host)")
	editCmd.Flags().StringVar(&editHostname, "hostname", "", "主机地址 (HostName)")
	editCmd.Flags().StringVar(&editUser, "user", "", "用户名 (User)")
	editCmd.Flags().StringVar(&editPort, "port", "", "端口 (Port)")
}
