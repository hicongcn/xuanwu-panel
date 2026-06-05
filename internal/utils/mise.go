package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var nodePathCache sync.Map
var misePath string
var miseOnce sync.Once

func FindMise() string {
	miseOnce.Do(func() {
		if p, err := exec.LookPath("mise"); err == nil {
			misePath = p
			return
		}
		home, _ := os.UserHomeDir()
		suffix := ""
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}
		candidates := []string{
			filepath.Join(home, ".local", "bin", "mise"+suffix),
			"/opt/homebrew/bin/mise",
			"/usr/local/bin/mise",
		}
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && info.Mode()&0111 != 0 {
				misePath = c
				return
			}
		}
		misePath = "mise" + suffix
	})
	return misePath
}

func MiseEnv() []string {
	env := os.Environ()
	home, _ := os.UserHomeDir()
	extra := filepath.Join(home, ".local", "bin")
	pathSep := ":"
	if runtime.GOOS == "windows" {
		pathSep = ";"
	}
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			if !strings.Contains(e, extra) {
				env[i] = "PATH=" + extra + pathSep + strings.TrimPrefix(e, "PATH=")
			}
			break
		}
	}
	if _, ok := lookupEnv(env, "MISE_CONFIG_DIR"); !ok {
		if !isDirWritable(filepath.Join(home, ".config", "mise")) {
			configDir := filepath.Join(home, ".local", "share", "mise")
			env = append(env, "MISE_CONFIG_DIR="+configDir)
		}
	}
	return env
}

func isDirWritable(dir string) bool {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false
	}
	f := filepath.Join(dir, ".write_test")
	if err := os.WriteFile(f, []byte("ok"), 0644); err != nil {
		return false
	}
	os.Remove(f)
	return true
}

func lookupEnv(env []string, key string) (string, bool) {
	prefix := key + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return strings.TrimPrefix(e, prefix), true
		}
	}
	return "", false
}

// GetMiseNodePath 获取指定版本的 node 全局包路径，使用内存缓存避免重复获取
func GetMiseNodePath(version string) string {
	if version == "" {
		version = "latest"
	}

	if val, ok := nodePathCache.Load(version); ok {
		return val.(string)
	}

	cmd := exec.Command("mise", "where", "node@"+version)
	out, err := cmd.CombinedOutput()
	if err == nil {
		nodeDir := strings.TrimSpace(string(out))
		if nodeDir != "" {
			// 采用双路径策略：lib/node_modules 是标准路径，lib 是某些环境（如 mise Docker）下的特殊路径
			// 通过冒号分隔，让 Node.js 按顺序搜索，保证最大兼容性
			nodePath := nodeDir + "/lib/node_modules:" + nodeDir + "/lib"
			nodePathCache.Store(version, nodePath)
			return nodePath
		}
	}

	return ""
}

// InjectNodePath 检查语言环境中是否有 node，如果有则自动获取并注入 NODE_PATH 到环境变量切片中
func InjectNodePath(envs *[]string, languages []map[string]string) {
	for _, lang := range languages {
		if lang["name"] == "node" {
			if nodePath := GetMiseNodePath(lang["version"]); nodePath != "" {
				*envs = append(*envs, "NODE_PATH="+nodePath)
			}
			break
		}
	}
}

// BuildMiseCommand 构建多语言 mise 执行命令 (字符串形式)
func BuildMiseCommand(command string, languages []map[string]string) string {
	if len(languages) == 0 {
		return command
	}

	var builder strings.Builder
	builder.WriteString("mise exec")

	for _, lang := range languages {
		name := lang["name"]
		version := lang["version"]
		if name == "" {
			continue
		}
		if version == "" {
			version = "latest"
		}
		builder.WriteString(" " + name + "@" + version)
	}

	builder.WriteString(" -- " + command)
	return builder.String()
}

// BuildMiseCommandSimple 构建单个语言的 mise 执行命令
func BuildMiseCommandSimple(command string, language, version string) string {
	if language == "" {
		return command
	}
	spec := language
	if version != "" {
		spec += "@" + version
	}
	return "mise exec " + spec + " -- " + command
}

// BuildMiseCommandArgsSimple 构建单个语言的 mise 执行命令 (参数列表形式)
func BuildMiseCommandArgsSimple(cmdArgs []string, language, version string) []string {
	if language == "" {
		return cmdArgs
	}
	spec := language
	if version != "" {
		spec += "@" + version
	}
	return append([]string{"mise", "exec", spec, "--"}, cmdArgs...)
}

// ListMiseInstalledVersions 获取指定语言已安装的所有版本列表
func ListMiseInstalledVersions(language string) ([]string, error) {
	// 执行 mise ls <language> 命令
	cmd := exec.Command(FindMise(), "ls", language)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var versions []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		v := strings.TrimSpace(line)
		if v == "" {
			continue
		}
		// mise ls 的输出可能包含状态标识或插件名，例:
		// * 20.10.0 (active)
		// node 18.17.0
		fields := strings.Fields(v)
		if len(fields) == 0 {
			continue
		}

		startIdx := 0
		// 跳过状态标识符
		if fields[startIdx] == "*" || fields[startIdx] == "->" || fields[startIdx] == ">" {
			startIdx++
		}

		if len(fields) <= startIdx {
			continue
		}

		vstr := fields[startIdx]
		// 如果第一个有效字段是插件名，则版本号在第二个字段
		if vstr == language && len(fields) > startIdx+1 {
			vstr = fields[startIdx+1]
		}

		// 确保解析出来的不是插件名
		if vstr != "" && vstr != language {
			versions = append(versions, vstr)
		}
	}
	return versions, nil
}
