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
	pushTags      string
	pushOr        bool
	pushParallel  int
	pushRetry     int
	pushRecursive bool
)

var pushCmd = &cobra.Command{
	Use:   "push [flags] <src> <dest>",
	Short: "批量上传文件/目录到主机 (scp)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]
		hostObjs, err := config.ParseSSHConfig()
		if err != nil {
			fmt.Println("解析 ssh config 失败:", err)
			os.Exit(1)
		}
		var filterTags []string
		if pushTags != "" {
			filterTags = strings.Split(pushTags, ",")
			for i := range filterTags {
				filterTags[i] = strings.TrimSpace(filterTags[i])
			}
		}
		var hosts []string
		if len(filterTags) > 0 {
			nameList, _ := tags.HostsByTags(filterTags, pushOr)
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
		results := runner.RunScpOnHosts(hosts, src, dest, runner.ScpPush, pushRecursive, pushParallel, pushRetry)
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("%s 失败: %v\n", r.Host, r.Err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&pushTags, "tag", "t", "", "按标签筛选主机，逗号分隔")
	pushCmd.Flags().BoolVar(&pushOr, "or", false, "标签筛选使用 OR 逻辑（默认 AND）")
	pushCmd.Flags().IntVar(&pushParallel, "parallel", 10, "并发数")
	pushCmd.Flags().IntVar(&pushRetry, "retry", 3, "失败重试次数")
	pushCmd.Flags().BoolVarP(&pushRecursive, "recursive", "r", false, "递归目录")
}
