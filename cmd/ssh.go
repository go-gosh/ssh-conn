package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/go-gosh/ssh-conn/internal/config"
	"github.com/go-gosh/ssh-conn/internal/picker"
	"github.com/go-gosh/ssh-conn/internal/tags"
	"github.com/spf13/cobra"
)

var (
	sshTags string
	sshOr   bool
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "交互式选择主机并登录",
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := config.ParseSSHConfig()
		if err != nil {
			fmt.Println("解析 ssh config 失败:", err)
			os.Exit(1)
		}
		var filterTags []string
		if sshTags != "" {
			filterTags = strings.Split(sshTags, ",")
			for i := range filterTags {
				filterTags[i] = strings.TrimSpace(filterTags[i])
			}
		}
		var filtered []config.Host
		if len(filterTags) > 0 {
			// 标签过滤
			nameList, _ := tags.HostsByTags(filterTags, sshOr)
			nameSet := map[string]struct{}{}
			for _, n := range nameList {
				nameSet[n] = struct{}{}
			}
			for _, h := range hosts {
				if _, ok := nameSet[h.Name]; ok {
					filtered = append(filtered, h)
				}
			}
		} else {
			filtered = hosts
		}
		if len(filtered) == 0 {
			fmt.Println("无匹配主机")
			os.Exit(1)
		}
		// 交互选择
		idx, err := picker.PickHost(filtered)
		if err != nil {
			fmt.Println("选择主机失败:", err)
			os.Exit(1)
		}
		h := filtered[idx]
		sshArgs := []string{h.Name}
		sshPath, err := exec.LookPath("ssh")
		if err != nil {
			fmt.Println("未找到 ssh 命令")
			os.Exit(1)
		}
		syscall.Exec(sshPath, append([]string{"ssh"}, sshArgs...), os.Environ())
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&sshTags, "tag", "t", "", "按标签筛选主机，逗号分隔")
	sshCmd.Flags().BoolVar(&sshOr, "or", false, "标签筛选使用 OR 逻辑（默认 AND）")
}
