# ssh-conn

一个用 Go 实现的命令行 SSH 主机管理工具，支持基于 `~/.ssh/config` 的主机配置管理和标签分组。

## 功能特性
- 列出所有主机，支持按标签筛选
- 新增、编辑、删除主机配置
- 多标签分组管理（基于 sqlite，纯 Go 实现，无需 cgo）
- 通过标签批量筛选主机
- 一键 ssh 连接主机
- 标签的增删改查
- 单元测试覆盖核心功能

## 安装

```bash
git clone https://github.com/go-gosh/ssh-conn.git
cd ssh-conn
go build -o ssh-conn
```

## 使用示例

### 主机管理

```bash
# 列出所有主机
./ssh-conn list

# 按标签筛选
./ssh-conn list --tag web

# 添加主机
./ssh-conn add --host myserver --hostname 1.2.3.4 --user root --tag web

# 编辑主机
./ssh-conn edit --host myserver --user admin

# 删除主机
./ssh-conn delete --host myserver
```

### 标签管理

```bash
# 添加标签（支持多个标签）
./ssh-conn tag add myserver --tags dev,web

# 移除标签
./ssh-conn tag remove myserver --tags web

# 覆盖标签
./ssh-conn tag edit myserver --tags prod

# 查询主机标签
./ssh-conn tag list myserver
```

### SSH 连接

```bash
./ssh-conn connect myserver
```

## 运行测试

```bash
go test ./db/...
```

## 依赖
- Go 1.21+
- [spf13/cobra](https://github.com/spf13/cobra)
- [gorm](https://gorm.io/)
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
- [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config)

## License
MIT 