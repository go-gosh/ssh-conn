# ssh-conn

一个基于 Go 的 SSH 主机批量管理与标签工具。

- 100% 复用 ~/.ssh/config 主机配置
- 标签独立存储于 ~/.ssh-conn/tags.yaml
- 支持标签批量管理、主机筛选、批量命令、文件分发等

## 依赖
- Go 1.22+
- cobra
- github.com/kevinburke/ssh_config
- gopkg.in/yaml.v3

## 目录结构
```
ssh-conn/
 ├─ cmd/               # cobra 子命令
 ├─ internal/
 │   ├─ config/        # 解析 ssh config
 │   ├─ tags/          # 读写 tags.yaml
 │   ├─ runner/        # 并发 ssh/scp
 │   ├─ picker/        # fzf/内置选择器
 │   └─ ui/            # 颜色/进度条
 ├─ go.mod
 ├─ Makefile
 └─ README.md
```

## 快速开始
```bash
make build
./ssh-conn list -t web,prod
./ssh-conn tag add prod web01 web02
./ssh-conn tag rm prod web01
./ssh-conn tag ls
```

## 标签文件示例
```yaml
web:  [web01, web02, web03]
db:   [db01, db02]
prod: [web01, db01]
```

## 贡献
欢迎 PR 和 issue！ 