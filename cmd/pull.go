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
	pullTags      string
	pullOr        bool
	pullParallel  int
	pullRetry     int
	pullRecursive bool
)

var pullCmd = &cobra.Command{
	Use:   "pull [flags] <src> <dest>",
	Short: "批量下载文件/目录到本地 (scp)",
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
		if pullTags != "" {
			filterTags = strings.Split(pullTags, ",")
			for i := range filterTags {
				filterTags[i] = strings.TrimSpace(filterTags[i])
			}
		}
		var hosts []string
		if len(filterTags) > 0 {
			nameList, _ := tags.HostsByTags(filterTags, pullOr)
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
		results := runner.RunScpOnHosts(hosts, src, dest, runner.ScpPull, pullRecursive, pullParallel, pullRetry)
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("%s 失败: %v\n", r.Host, r.Err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&pullTags, "tag", "t", "", "按标签筛选主机，逗号分隔")
	pullCmd.Flags().BoolVar(&pullOr, "or", false, "标签筛选使用 OR 逻辑（默认 AND）")
	pullCmd.Flags().IntVar(&pullParallel, "parallel", 10, "并发数")
	pullCmd.Flags().IntVar(&pullRetry, "retry", 3, "失败重试次数")
	pullCmd.Flags().BoolVarP(&pullRecursive, "recursive", "r", false, "递归目录")
}
