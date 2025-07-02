package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var delHost string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除 ssh config 中的主机配置",
	Run: func(cmd *cobra.Command, args []string) {
		if delHost == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "--host 必填")
			return
		}
		configPath := os.ExpandEnv("$HOME/.ssh/config")
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
			if strings.HasPrefix(trim, "Host ") && strings.TrimSpace(trim[5:]) == delHost {
				inBlock = true
				continue
			}
			if inBlock && (strings.HasPrefix(trim, "Host ") && strings.TrimSpace(trim[5:]) != delHost) {
				inBlock = false
			}
			if !inBlock {
				out = append(out, line)
			}
		}
		os.WriteFile(configPath, []byte(strings.Join(out, "\n")), 0600)
		fmt.Fprintln(cmd.OutOrStdout(), "删除完成")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&delHost, "host", "", "主机别名 (Host)")
}
