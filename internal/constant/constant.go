package constant

import "time"

const (

	// DefaultRole 默认用户角色
	DefaultRole = "user"

	// AdminRole 管理员角色
	AdminRole = "admin"

	// CookieName Cookie 名称
	CookieName = "XWToken"

	// DefaultTaskTimeout 默认任务超时时间（分钟）
	DefaultTaskTimeout = 30

	// Settings Section 常量
	SectionSite      = "site"
	SectionSystem    = "system"
	SectionScheduler = "scheduler"
	SectionSecurity  = "security"
	SectionNotify    = "notify"

	// Site Settings Key 常量
	KeyTitle        = "title"
	KeySubtitle     = "subtitle"
	KeyIcon         = "icon"
	KeyPageSize     = "page_size"
	KeyCookieDays   = "cookie_days"
	KeyOpenapiToken = "openapi_token"

	// Security Settings Key 常量
	KeySecret = "secret"

	// System Settings Key 常量
	KeyInitialized = "initialized"
	// Log Retention Keys
	KeySystemNoticeDays     = "system_notice_days"
	KeySystemNoticeMaxCount = "system_notice_max_count"
	KeyPushLogDays          = "push_log_days"
	KeyPushLogMaxCount      = "push_log_max_count"
	KeyLoginLogDays         = "login_log_days"
	KeyLoginLogMaxCount     = "login_log_max_count"
	KeySchedulerLogDays     = "scheduler_log_days"
	KeySchedulerLogMaxCount = "scheduler_log_max_count"

	// Scheduler Settings Key 常量
	KeyWorkerCount  = "worker_count"
	KeyQueueSize    = "queue_size"
	KeyRateInterval = "rate_interval"

	// Notify Settings Key 常量
	KeyNotifyChannels = "channels"
	KeyNotifyEvents   = "events"
	KeyNotifyToken    = "notify_token"
	KeyNotifyPrefix   = "notify_prefix"

	// Notify Templates Keys
	KeyNotifyTemplateUserLoginTitle       = "notify_template_user_login_title"
	KeyNotifyTemplateUserLoginText        = "notify_template_user_login_text"
	KeyNotifyTemplateBruteForceLoginTitle = "notify_template_brute_force_login_title"
	KeyNotifyTemplateBruteForceLoginText  = "notify_template_brute_force_login_text"
	KeyNotifyTemplatePasswordChangedTitle = "notify_template_password_changed_title"
	KeyNotifyTemplatePasswordChangedText  = "notify_template_password_changed_text"
	KeyNotifyTemplateTaskSuccessTitle     = "notify_template_task_success_title"
	KeyNotifyTemplateTaskSuccessText      = "notify_template_task_success_text"
	KeyNotifyTemplateTaskFailedTitle      = "notify_template_task_failed_title"
	KeyNotifyTemplateTaskFailedText       = "notify_template_task_failed_text"
	KeyNotifyTemplateTaskTimeoutTitle     = "notify_template_task_timeout_title"
	KeyNotifyTemplateTaskTimeoutText      = "notify_template_task_timeout_text"

	// 事件绑定类型
	BindingTypeSystem = "system"
	BindingTypeTask   = "task"

	// 系统事件类型
	EventUserLogin       = "user_login"
	EventBruteForceLogin = "brute_force_login"
	EventPasswordChanged = "password_changed"

	// 任务事件类型
	EventTaskSuccess = "task_success"
	EventTaskFailed  = "task_failed"
	EventTaskTimeout = "task_timeout"
	EventTaskRunning = "task_running"
	EventTaskQueued  = "task_queued"

	// 其他事件类型
	EventSystemNotice = "system_notice"
	EventSchedulerLog = "scheduler_log"
	EventNotifySent   = "notify_sent"
	EventAppLogAdded  = "app_log_added"

	// 任务状态
	TaskStatusSuccess   = "success"
	TaskStatusFailed    = "failed"
	TaskStatusRunning   = "running"
	TaskStatusPending   = "pending"
	TaskStatusTimeout   = "timeout"
	TaskStatusCancelled = "cancelled"
	TaskStatusQueued    = "queued"

	// 任务类型
	TaskTypeNormal = "task"
	TaskTypeRepo   = "repo"

	// 任务置顶类型
	PinTypeNone = "none"
	PinTypeTop  = "top"

	// 触发类型
	TriggerTypeCron          = "cron"
	TriggerTypeXuanwuStartup = "xuanwu_startup"

	// AppLog 分类
	LogCategoryDefault      = "default"
	LogCategorySystemNotice = "system_notice"
	LogCategoryPushLog      = "push_log"
	LogCategoryLoginLog     = "login_log"
	LogCategorySchedulerLog = "scheduler_log"

	// AppLog 级别
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"

	// AppLog 状态
	LogStatusUnread  = "unread"
	LogStatusRead    = "read"
	LogStatusSuccess = "success"
	LogStatusFailed  = "failed"

	// WebSocket 安全常量
	// PongWait 收到 pong 的超时时间
	PongWait = 60 * time.Second
	// PingPeriod 发送 ping 的周期
	PingPeriod = (PongWait * 9) / 10
	// MaxMessageSize 允许的最大消息大小
	MaxMessageSize = 1024 * 1024 // 1MB
	// MaxLogSize 允许的最大日志大小 (保留末尾 10MB)
	MaxLogSize = 10 * 1024 * 1024 // 10MB

	// ScriptsDirPlaceholder 脚本目录占位符
	ScriptsDirPlaceholder = "$SCRIPTS_DIR$"
)

// TablePrefix 表前缀，从配置文件读取
var TablePrefix string

// Runtime 数据库配置快照，用于需要单独启动内部子进程（如 reposync）时显式透传数据库连接信息，
// 避免主进程启动阶段清理环境变量后，子进程意外回退到默认 sqlite 配置。
var (
	RuntimeDBType        string
	RuntimeDBHost        string
	RuntimeDBPort        int
	RuntimeDBUser        string
	RuntimeDBPassword    string
	RuntimeDBName        string
	RuntimeDBPath        string
	RuntimeDBDSN         string
	RuntimeDBTablePrefix string
)

// Secret JWT和密码salt密钥，运行中自动从数据库加载
var Secret string

// DemoMode 演示模式，从环境变量读取
var DemoMode bool

// DefaultIcon 默认站点图标
var DefaultIcon = `<svg height="800px" width="800px" version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512" xml:space="preserve"><g><path style="fill:#76CFE2;" d="M10.199,288.583c0-42.869,34.752-77.621,77.621-77.621s77.621,34.752,77.621,77.621H10.199z"/><path style="fill:#76CFE2;" d="M346.558,288.583c0-42.869,34.752-77.621,77.621-77.621s77.622,34.753,77.622,77.621H346.558z"/><path style="fill:#76CFE2;" d="M259.105,389.975c30.313-30.313,79.46-30.313,109.773,0c30.313,30.313,30.313,79.46,0,109.773L259.105,389.975z"/><path style="fill:#76CFE2;" d="M254.999,389.975c-30.313-30.313-79.46-30.313-109.773,0c-30.313,30.313-30.313,79.46,0,109.773L254.999,389.975z"/><path style="fill:#76CFE2;" d="M307.748,174.739H204.252V64c0-28.579,23.169-51.748,51.748-51.748l0,0c28.579,0,51.748,23.169,51.748,51.748V174.739z"/></g><circle style="fill:#FF931E;" cx="256" cy="277.204" r="170.245"/><polygon style="fill:#F15A24;" points="174.245,229.998 256,182.797 337.755,229.998 337.755,324.4 256,371.602 174.245,324.4 "/><path d="M424.18,200.763c-1.558,0-3.124,0.043-4.685,0.125c-20.129-42.954-56.688-76.736-101.546-93.182V64c0-34.157-27.79-61.947-61.947-61.947S194.054,29.843,194.054,64v43.707c-44.86,16.447-81.419,50.231-101.548,93.186c-1.56-0.083-3.125-0.129-4.684-0.129C39.396,200.763,0,240.16,0,288.583c0,5.633,4.566,10.199,10.199,10.199h66.648c4.678,39.071,21.889,74.34,47.486,101.669c-9.018,15.387-13.176,33.319-11.773,51.284c1.626,20.824,10.666,40.436,25.454,55.225c1.992,1.992,4.602,2.987,7.212,2.987s5.221-0.995,7.212-2.987l55.753-55.753c15.234,4.19,31.263,6.441,47.809,6.441c17.145,0,33.738-2.405,49.461-6.893l56.205,56.205c1.992,1.992,4.602,2.987,7.212,2.987s5.221-0.995,7.212-2.987c14.789-14.789,23.828-34.403,25.454-55.229c1.444-18.5-3.024-36.957-12.609-52.643c24.899-27.134,41.617-61.883,46.218-100.306h66.649c5.633,0,10.199-4.566,10.199-10.199C512,240.16,472.604,200.763,424.18,200.763z M214.452,64c0-22.909,18.638-41.548,41.548-41.548c22.91,0,41.548,18.638,41.548,41.548v37.587c-13.345-3.157-27.253-4.835-41.548-4.835c-14.295,0-28.204,1.679-41.548,4.835V64z M21.168,278.384c4.765-31.266,31.086-55.537,63.273-57.134c-5.762,17.623-8.891,36.426-8.891,55.949c0,0.397,0.012,0.789,0.015,1.185L21.168,278.384z M121.618,190.36l22.637,22.637c1.992,1.992,4.602,2.987,7.212,2.987c2.61,0,5.221-0.995,7.212-2.987c3.983-3.983,3.983-10.441,0-14.425l-24.756-24.756c29.381-34.641,73.201-56.669,122.075-56.669s92.694,22.028,122.075,56.669l-42.909,42.909l-68.966-39.817v-11.484c0-5.633-4.566-10.199-10.199-10.199s-10.199,4.566-10.199,10.199v11.484l-76.655,44.257c-3.156,1.822-5.1,5.188-5.1,8.833v88.508l-43.635,43.635c-15.487-24.634-24.462-53.757-24.462-84.942C95.95,245.211,105.388,215.39,121.618,190.36z M267.234,436.849v-59.957l69.244-39.979l42.545,42.545C351.888,412.049,312.077,433.724,267.234,436.849z M132.503,378.897l42.362-42.362l71.97,41.552v58.893C200.889,434.374,160.062,412.305,132.503,378.897z M256,194.574l71.556,41.313v82.625L256,359.825l-71.556-41.313v-82.625L256,194.574z M145.763,484.786c-15.009-20.324-17.359-47.248-6.136-69.797c14.036,11.873,29.899,21.646,47.107,28.826L145.763,484.786z M368.341,484.788l-41.594-41.594c17.233-7.373,33.095-17.35,47.087-29.439C385.773,436.55,383.627,464.09,368.341,484.788z M391.19,362.777l-43.236-43.236v-86.753l42.428-42.428c16.23,25.029,25.668,54.851,25.668,86.839C416.05,308.652,406.925,338.011,391.19,362.777z M436.434,278.384c0.003-0.396,0.015-0.789,0.015-1.185c0-19.522-3.129-38.326-8.891-55.949c32.187,1.597,58.508,25.869,63.273,57.134H436.434z"/><path d="M256,147.68c5.633,0,10.199-4.566,10.199-10.199v-2.07c0-5.633-4.566-10.199-10.199-10.199s-10.199,4.566-10.199,10.199v2.07C245.801,143.114,250.367,147.68,256,147.68z"/></svg>`

// DefaultSettings 默认系统设置
var DefaultSettings = map[string]map[string]string{
	SectionSite: {
		KeyTitle:      "玄武面板",
		KeySubtitle:   "极致轻量、高性能的自动化任务调度平台",
		KeyIcon:       DefaultIcon,
		KeyPageSize:   "10",
		KeyCookieDays: "7",
	},
	SectionScheduler: {
		KeyWorkerCount:  "4",
		KeyQueueSize:    "100",
		KeyRateInterval: "200",
	},
	SectionNotify: {
		KeyNotifyPrefix: "[玄武面板]",
		// Login
		KeyNotifyTemplateUserLoginTitle:       "用户登录(成功/失败)",
		KeyNotifyTemplateUserLoginText:        "用户 {{username}} 在 IP {{ip}} 登录{{status_label}}\n{{message}}",
		KeyNotifyTemplateBruteForceLoginTitle: "系统安全警告",
		KeyNotifyTemplateBruteForceLoginText:  "检测到 IP {{ip}} 正在尝试暴力破解用户 {{username}}",
		KeyNotifyTemplatePasswordChangedTitle: "账户安全通知",
		KeyNotifyTemplatePasswordChangedText:  "用户 {{username}} 刚刚修改了密码",
		// Task
		KeyNotifyTemplateTaskSuccessTitle: "任务[{{task_name}}] 成功",
		KeyNotifyTemplateTaskSuccessText:  "任务 #{{task_id}} {{task_name}}\n状态: 成功\n耗时: {{duration}}ms\n执行结果: {{output}}",
		KeyNotifyTemplateTaskFailedTitle:  "任务[{{task_name}}] 失败",
		KeyNotifyTemplateTaskFailedText:   "任务 #{{task_id}} {{task_name}}\n状态: 失败\n执行时间: {{start_time}}\n原因: {{error}}\n最后输出: {{output}}",
		KeyNotifyTemplateTaskTimeoutTitle: "任务[{{task_name}}] 超时",
		KeyNotifyTemplateTaskTimeoutText:  "任务 #{{task_id}} {{task_name}}\n状态: 超时\n耗时: {{duration}}ms\n最后输出: {{output}}",
	},
}
