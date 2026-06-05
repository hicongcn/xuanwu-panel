package utils

// MaskString 对字符串进行脱敏处理，保留首尾，中间用星号遮掩
func MaskString(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	if n <= 3 {
		return "***"
	}
	if n <= 6 {
		return s[:1] + "***" + s[n-1:]
	}
	return s[:2] + "****" + s[n-2:]
}
