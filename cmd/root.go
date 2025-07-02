package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshconfig",
	Short: "SSH 主机配置与标签管理工具",
	Long:  `一个用于管理 ~/.ssh/config 和主机标签的命令行工具。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
