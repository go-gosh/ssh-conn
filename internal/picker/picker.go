package picker

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/go-gosh/ssh-conn/internal/config"
	"github.com/ktr0731/go-fuzzyfinder"
)

// PickHost 交互选择主机，返回所选索引
func PickHost(hosts []config.Host) (int, error) {
	// 优先用 fzf
	if _, err := exec.LookPath("fzf"); err == nil {
		cmd := exec.Command("fzf")
		in, _ := cmd.StdinPipe()
		out, _ := cmd.StdoutPipe()
		cmd.Start()
		for _, h := range hosts {
			in.Write([]byte(h.Name + "\n"))
		}
		in.Close()
		var sel string
		fmt.Fscanln(out, &sel)
		cmd.Wait()
		for i, h := range hosts {
			if h.Name == sel {
				return i, nil
			}
		}
		return -1, errors.New("fzf 未选中主机")
	}
	// 否则用 go-fuzzyfinder
	idx, err := fuzzyfinder.Find(
		hosts,
		func(i int) string { return hosts[i].Name },
		fuzzyfinder.WithPromptString("选择主机: "),
	)
	return idx, err
}
