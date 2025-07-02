package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var sshExec = exec.Command

var connectCmd = &cobra.Command{
	Use:   "connect [host]",
	Short: "通过 ssh 连接指定主机",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		sshPath, err := exec.LookPath("ssh")
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "未找到 ssh 命令")
			return
		}
		c := sshExec(sshPath, host)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "ssh 连接失败:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
