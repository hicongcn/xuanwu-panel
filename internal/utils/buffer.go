package utils

import "fmt"

// TrimLog 裁剪日志，保留末尾指定大小
func TrimLog(content string, limit int) string {
	if len(content) <= limit {
		return content
	}
	// 简单裁剪，不考虑字符完整性，因为这是针对大文本的保护
	return fmt.Sprintf("\n\n[System] 日志过长，已自动截断，仅保留末尾 %d MB...\n\n", limit/1024/1024) + content[len(content)-limit:]
}
