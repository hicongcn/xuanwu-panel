//go:build !web

package static

import (
	"io/fs"
	"os"
)

// DefaultStaticDir 是 Docker 模式下默认的静态资源目录
const DefaultStaticDir = "/www/xuanwu"

func GetFS() fs.FS {
	dir := getStaticDir()
	if dir == "" {
		return nil
	}
	return os.DirFS(dir)
}

func getStaticDir() string {
	// 优先从环境变量读取
	if dir := os.Getenv("XW_STATIC_DIR"); dir != "" {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	// 回退到默认目录
	if _, err := os.Stat(DefaultStaticDir); err == nil {
		return DefaultStaticDir
	}
	return ""
}
