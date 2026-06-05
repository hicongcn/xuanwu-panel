package constant

import "time"

// 构建时注入的变量
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// 程序启动时间
var StartTime = time.Now()
