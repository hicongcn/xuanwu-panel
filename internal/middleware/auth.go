package middleware

import (
	"net/http"

	"github.com/hicongcn/xuanwu-panel/internal/constant"
	"github.com/hicongcn/xuanwu-panel/internal/database"
	"github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 基础的 CSRF 防护：校验 Origin/Referer (针对非 GET 请求)
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodOptions && c.Request.Method != http.MethodHead {
			origin := c.GetHeader("Origin")
			if origin == "" {
				origin = c.GetHeader("Referer")
			}
			// 如果有 Origin 且不匹配则拒绝（实际部署时应配置允许的 Origin）
			if origin != "" && !utils.CheckWSOrigin(c.Request) {
				utils.Forbidden(c, "CSRF 校验失败: 非法的请求来源")
				c.Abort()
				return
			}
		}

		token, err := c.Cookie(constant.CookieName)
		if err != nil || token == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 验证 token
		userID, username, tokenVersion, err := utils.ParseToken(token, constant.Secret)
		if err != nil {
			utils.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// 安全增强：校验数据库中该用户的 ID 是否与 Token 一致，并验证 TokenVersion
		var user models.User
		res := database.DB.Where("username = ?", username).Limit(1).Find(&user)
		if res.Error != nil || res.RowsAffected == 0 || user.ID != userID || user.TokenVersion != tokenVersion {
			utils.Unauthorized(c, "会话失效，请重新登录")
			ClearAuthCookie(c)
			c.Abort()
			return
		}

		// 将用户信息存入上下文 (必须使用数据库中的最新 ID)
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)
		c.Next()
	}
}

// AdminRequired 管理员权限认证中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != constant.AdminRole {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetAuthCookie 设置认证 Cookie，expireDays 为过期天数
func SetAuthCookie(c *gin.Context, token string, expireDays int) {
	maxAge := 86400 * expireDays
	// 增加 SameSite=Lax 和 Secure 属性（如果环境支持，这里暂时设为 false，但生产建议 true）
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(constant.CookieName, token, maxAge, "/", "", false, true)
}

// ClearAuthCookie 清除认证 Cookie
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie(constant.CookieName, "", -1, "/", "", false, true)
}

// LocalhostOnly 仅允许本地回环地址访问，并进行简单的内部凭证校验
func LocalhostOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip != "127.0.0.1" && ip != "::1" {
			utils.BadRequest(c, "仅允许本地访问")
			c.Abort()
			return
		}

		// 简单的内部通信认证
		token := c.GetHeader("X-Internal-Token")
		if token == "" || token != constant.Secret {
			utils.Unauthorized(c, "无效的内部调用凭证")
			c.Abort()
			return
		}

		c.Next()
	}
}
