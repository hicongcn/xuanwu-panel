package vo

import (
	"github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/utils"
)

// UserVO 用户视图对象
type UserVO struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	CreatedAt models.LocalTime `json:"created_at"`
	UpdatedAt models.LocalTime `json:"updated_at"`
}

// EnvVO 环境变量视图对象
type EnvVO struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Value     string           `json:"value"`
	Remark    string           `json:"remark"`
	Type      string           `json:"type"`
	Hidden    bool             `json:"hidden"`
	Enabled   bool             `json:"enabled"`
	CreatedAt models.LocalTime `json:"created_at"`
	UpdatedAt models.LocalTime `json:"updated_at"`
}

// ToEnvVO 将 Env 模型转换为 EnvVO
func ToEnvVO(env *models.EnvironmentVariable) *EnvVO {
	if env == nil {
		return nil
	}
	val := string(env.Value)
	return &EnvVO{
		ID:        env.ID,
		Name:      env.Name,
		Value:     val,
		Remark:    env.Remark,
		Type:      env.Type,
		Hidden:    utils.DerefBool(env.Hidden, true),
		Enabled:   utils.DerefBool(env.Enabled, true),
		CreatedAt: env.CreatedAt,
		UpdatedAt: env.UpdatedAt,
	}
}

// ToEnvVOListFromModels 将 Env 模型列表转换为 EnvVO 列表
func ToEnvVOListFromModels(envs []models.EnvironmentVariable) []*EnvVO {
	vos := make([]*EnvVO, len(envs))
	for i := range envs {
		vos[i] = ToEnvVO(&envs[i])
	}
	return vos
}

// LoginLogVO 登录日志视图对象
type LoginLogVO struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	IP        string           `json:"ip"`
	UserAgent string           `json:"user_agent"`
	Status    string           `json:"status"`
	Message   string           `json:"message"`
	CreatedAt models.LocalTime `json:"created_at"`
}

// TokenConfig Token 配置结构体
type TokenConfig struct {
	Enabled  bool   `json:"enabled"`
	Token    string `json:"token"`
	ExpireAt string `json:"expire_at"`
}
