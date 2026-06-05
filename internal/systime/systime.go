package systime

import "time"

// CST 东八区时区
var CST = time.FixedZone("CST", 8*3600)

// Now 返回东八区的当前时间
func Now() time.Time {
	return time.Now().In(CST)
}

// InCST 将给定时间转换为东八区时间
func InCST(t time.Time) time.Time {
	return t.In(CST)
}

// FormatDate 将时间格式化为日期格式
func FormatDate(t time.Time) string {
	return InCST(t).Format("2006-01-02")
}

// FormatDatetime 用于备份等文件命名格式
func FormatDatetime(t time.Time) string {
	return InCST(t).Format("20060102_150405")
}
