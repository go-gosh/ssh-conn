package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type OutputMode int

const (
	OutputColor OutputMode = iota
	OutputJSON
)

type RunResult struct {
	Host   string `json:"host"`
	Line   string `json:"line"`
	Err    error  `json:"err,omitempty"`
	Stdout bool   `json:"stdout"`
}

// RunOnHosts 并发在 hosts 上执行 cmdStr，支持重试和实时输出
func RunOnHosts(hosts []string, userCmd string, parallel, retry int, mode OutputMode) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, host := range hosts {
		host := host
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var lastErr error
			for i := 0; i < retry; i++ {
				cmd := exec.CommandContext(ctx, "ssh", host, userCmd)
				stdout, _ := cmd.StdoutPipe()
				stderr, _ := cmd.StderrPipe()
				if err := cmd.Start(); err != nil {
					lastErr = err
					time.Sleep(500 * time.Millisecond)
					continue
				}
				outCh := make(chan RunResult)
				go func() {
					scan := bufio.NewScanner(stdout)
					for scan.Scan() {
						outCh <- RunResult{Host: host, Line: scan.Text(), Stdout: true}
					}
				}()
				go func() {
					scan := bufio.NewScanner(stderr)
					for scan.Scan() {
						outCh <- RunResult{Host: host, Line: scan.Text(), Stdout: false}
					}
				}()
				cmd.Wait()
				close(outCh)
				for r := range outCh {
					printResult(r, mode)
				}
				lastErr = nil
				break
			}
			if lastErr != nil {
				printResult(RunResult{Host: host, Line: lastErr.Error(), Err: lastErr}, mode)
			}
		}()
	}
	wg.Wait()
}

func printResult(r RunResult, mode OutputMode) {
	if mode == OutputJSON {
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
		return
	}
	color := "\033[36m"
	if !r.Stdout {
		color = "\033[31m"
	}
	fmt.Printf("%s%s\033[0m: %s\n", color, r.Host, r.Line)
}

// ScpMode 表示 push 还是 pull
type ScpMode int

const (
	ScpPush ScpMode = iota
	ScpPull
)

type ScpResult struct {
	Host string
	Err  error
}

// RunScpOnHosts 并发 push/pull 文件/目录，支持递归、重试、进度条
func RunScpOnHosts(hosts []string, src, dest string, mode ScpMode, recursive bool, parallel, retry int) []ScpResult {
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)
	results := make([]ScpResult, len(hosts))
	bar := progressbar.NewOptions(len(hosts), progressbar.OptionSetDescription("SCP"))
	for i, host := range hosts {
		i, host := i, host
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem; bar.Add(1) }()
			var lastErr error
			for t := 0; t < retry; t++ {
				var scpArgs []string
				if recursive {
					scpArgs = append(scpArgs, "-r")
				}
				if mode == ScpPush {
					scpArgs = append(scpArgs, src, fmt.Sprintf("%s:%s", host, dest))
				} else {
					scpArgs = append(scpArgs, fmt.Sprintf("%s:%s", host, src), dest)
				}
				cmd := exec.Command("scp", scpArgs...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					lastErr = err
					continue
				} else {
					lastErr = nil
					break
				}
			}
			results[i] = ScpResult{Host: host, Err: lastErr}
		}()
	}
	wg.Wait()
	return results
}
