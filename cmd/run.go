package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gosh/ssh-conn/internal/config"
	"github.com/go-gosh/ssh-conn/internal/runner"
	"github.com/go-gosh/ssh-conn/internal/tags"
	"github.com/spf13/cobra"
)

var (
	runTags     string
	runOr       bool
	runCmd      string
	runParallel int
	runRetry    int
	runJSON     bool
)

var runCmdCobra = &cobra.Command{
	Use:   "run",
	Short: "批量并发 ssh 执行命令",
	Run: func(cmd *cobra.Command, args []string) {
		if runCmd == "" {
			fmt.Println("请用 -c/--command 指定要执行的命令")
			os.Exit(1)
		}
		hostObjs, err := config.ParseSSHConfig()
		if err != nil {
			fmt.Println("解析 ssh config 失败:", err)
			os.Exit(1)
		}
		var filterTags []string
		if runTags != "" {
			filterTags = strings.Split(runTags, ",")
			for i := range filterTags {
				filterTags[i] = strings.TrimSpace(filterTags[i])
			}
		}
		var hosts []string
		if len(filterTags) > 0 {
			nameList, _ := tags.HostsByTags(filterTags, runOr)
			hostSet := map[string]struct{}{}
			for _, n := range nameList {
				hostSet[n] = struct{}{}
			}
			for _, h := range hostObjs {
				if _, ok := hostSet[h.Name]; ok {
					hosts = append(hosts, h.Name)
				}
			}
		} else {
			for _, h := range hostObjs {
				hosts = append(hosts, h.Name)
			}
		}
		if len(hosts) == 0 {
			fmt.Println("无匹配主机")
			os.Exit(1)
		}
		mode := runner.OutputColor
		if runJSON {
			mode = runner.OutputJSON
		}
		if runParallel <= 0 {
			runParallel = 10
		}
		if runRetry <= 0 {
			runRetry = 3
		}
		runner.RunOnHosts(hosts, runCmd, runParallel, runRetry, mode)
	},
}

func init() {
	rootCmd.AddCommand(runCmdCobra)
	runCmdCobra.Flags().StringVarP(&runTags, "tag", "t", "", "按标签筛选主机，逗号分隔")
	runCmdCobra.Flags().BoolVar(&runOr, "or", false, "标签筛选使用 OR 逻辑（默认 AND）")
	runCmdCobra.Flags().StringVarP(&runCmd, "command", "c", "", "要执行的命令")
	runCmdCobra.Flags().IntVar(&runParallel, "parallel", 10, "并发数")
	runCmdCobra.Flags().IntVar(&runRetry, "retry", 3, "失败重试次数")
	runCmdCobra.Flags().BoolVar(&runJSON, "json", false, "输出 JSON 格式")
}
