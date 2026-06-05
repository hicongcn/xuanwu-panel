package constant

import (
	"os"
	"path/filepath"
)

var (
	// ConfigPath 配置文件路径
	ConfigPath string

	// DataDir 数据目录
	DataDir string

	// DefaultDBPath 默认数据库路径
	DefaultDBPath string

	// WebDistDir 前端构建目录
	WebDistDir string

	// ScriptsWorkDir 脚本工作目录
	ScriptsWorkDir string

	// LogDir 日志目录
	LogDir string
)

func init() {
	rootDir := ResolveAppRootDir()
	DataDir = filepath.Clean(filepath.Join(rootDir, "data"))
	ConfigPath = filepath.Clean(filepath.Join(DataDir, "config", "config.ini"))
	DefaultDBPath = filepath.Clean(filepath.Join(DataDir, "db", "xuanwu.db"))
	WebDistDir = filepath.Clean(filepath.Join(rootDir, "web", "dist"))
	ScriptsWorkDir = filepath.Clean(filepath.Join(DataDir, "scripts"))
	LogDir = filepath.Clean(filepath.Join(DataDir, "log"))
}

// ResolveAppRootDir 获取应用程序的绝对根目录路径。
func ResolveAppRootDir() string {
	// 1. 检查当前工作目录（CWD）及其上级目录
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for {
			if _, err := os.Stat(filepath.Join(dir, "data", "config", "config.ini")); err == nil {
				return dir
			}
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// 2. 检查当前可执行文件路径及其上级目录
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		for {
			if _, err := os.Stat(filepath.Join(dir, "data", "config", "config.ini")); err == nil {
				return dir
			}
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// 3. 兜底回退到当前工作目录
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return "."
}
