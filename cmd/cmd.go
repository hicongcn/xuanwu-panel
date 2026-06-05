package cmd

import (
	"github.com/hicongcn/xuanwu-panel/cmd/builtininstall"
	"github.com/hicongcn/xuanwu-panel/cmd/reposync"
	"github.com/hicongcn/xuanwu-panel/cmd/resetpwd"
	"github.com/hicongcn/xuanwu-panel/cmd/restore"
	"github.com/hicongcn/xuanwu-panel/cmd/task"
	// "github.com/hicongcn/xuanwu-panel/cmd/migrate"
)

// CommandHandler 定义命令执行函数
type CommandHandler func(args []string)

// Handlers 维护了除了 server 之外的命令的执行入口
var Handlers = map[string]CommandHandler{
	"reposync":       reposync.Run,
	"resetpwd":       resetpwd.Run,
	"restore":        restore.Run,
	"builtininstall": builtininstall.Run,
	"task":           task.Run,
	// "migrate":  migrate.Run,
}
