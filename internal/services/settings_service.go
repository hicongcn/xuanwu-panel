package services

import (
	"github.com/hicongcn/xuanwu-panel/internal/cache"
	"github.com/hicongcn/xuanwu-panel/internal/constant"
	"github.com/hicongcn/xuanwu-panel/internal/database"
	"github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/utils"
)

type SettingsService struct{}

func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

// InitSettings 初始化默认设置
func (s *SettingsService) InitSettings() error {
	for section, keys := range constant.DefaultSettings {
		for key, value := range keys {
			var count int64
			database.DB.Model(&models.Setting{}).Where(&models.Setting{Section: section, Key: key}).Count(&count)
			if count == 0 {
				if err := database.DB.Create(&models.Setting{
					ID:      utils.GenerateID(),
					Section: section,
					Key:     key,
					Value:   models.BigText(value),
				}).Error; err != nil {
					return err
				}
			}
		}
	}

	// 默认值初始化
	defaultRetention := map[string]string{
		constant.KeySystemNoticeDays:     "30",
		constant.KeySystemNoticeMaxCount: "500",
		constant.KeyPushLogDays:          "15",
		constant.KeyPushLogMaxCount:      "5000",
		constant.KeyLoginLogDays:         "30",
		constant.KeyLoginLogMaxCount:     "1000",
		constant.KeySchedulerLogDays:     "30",
		constant.KeySchedulerLogMaxCount: "10000",
	}

	for k, v := range defaultRetention {
		var count int64
		database.DB.Model(&models.Setting{}).Where(&models.Setting{Section: constant.SectionSystem, Key: k}).Count(&count)
		if count == 0 {
			s.Set(constant.SectionSystem, k, v)
		}
	}

	// 从 constant.DefaultSettings 初始化所有缺少的通知模板
	if notifyDefaults, ok := constant.DefaultSettings[constant.SectionNotify]; ok {
		for k, v := range notifyDefaults {
			var count int64
			database.DB.Model(&models.Setting{}).Where(&models.Setting{Section: constant.SectionNotify, Key: k}).Count(&count)
			if count == 0 {
				s.Set(constant.SectionNotify, k, v)
			}
		}
	}

	// 初始化或获取 JWT Secret 密码
	var secCount int64
	database.DB.Model(&models.Setting{}).Where(&models.Setting{Section: constant.SectionSecurity, Key: constant.KeySecret}).Count(&secCount)
	var secretValue string
	if secCount == 0 {
		// 先尝试从配置文件读取遗留下来的旧设
		if Config != nil && Config.Security.Secret != "" {
			secretValue = Config.Security.Secret
		} else {
			secretValue = utils.RandomString(32)
		}
		if err := database.DB.Create(&models.Setting{
			ID:      utils.GenerateID(),
			Section: constant.SectionSecurity,
			Key:     constant.KeySecret,
			Value:   models.BigText(secretValue),
		}).Error; err != nil {
			return err
		}
	} else {
		secretValue = s.Get(constant.SectionSecurity, constant.KeySecret)
	}
	constant.Secret = secretValue

	cache.LoadSiteCache()
	return nil
}

// Get 获取单个设置
func (s *SettingsService) Get(section, key string) string {
	if section == constant.SectionSite {
		return cache.GetSiteCache(key)
	}
	var setting models.Setting
	res := database.DB.Where(&models.Setting{Section: section, Key: key}).Limit(1).Find(&setting)
	if res.Error != nil || res.RowsAffected == 0 {
		if def, ok := constant.DefaultSettings[section][key]; ok {
			return def
		}
		return ""
	}
	return string(setting.Value)
}

// Set 设置单个值
func (s *SettingsService) Set(section, key, value string) error {
	var setting models.Setting
	res := database.DB.Where(&models.Setting{Section: section, Key: key}).Limit(1).Find(&setting)
	var err error
	if res.Error != nil || res.RowsAffected == 0 {
		err = database.DB.Create(&models.Setting{
			ID:      utils.GenerateID(),
			Section: section,
			Key:     key,
			Value:   models.BigText(value),
		}).Error
	} else {
		err = database.DB.Model(&setting).Update("value", models.BigText(value)).Error
	}

	if err == nil && section == constant.SectionSite {
		cache.SetSiteCache(key, value)
	}
	return err
}

// Delete 删除单个设置
func (s *SettingsService) Delete(section, key string) error {
	return database.DB.Where(&models.Setting{Section: section, Key: key}).Delete(&models.Setting{}).Error
}

// GetSection 获取整个 section 的设置
func (s *SettingsService) GetSection(section string) map[string]string {
	if section == constant.SectionSite {
		return cache.GetSiteCacheAll()
	}
	result := make(map[string]string)
	if defaults, ok := constant.DefaultSettings[section]; ok {
		for k, v := range defaults {
			result[k] = v
		}
	}
	var settings []models.Setting
	database.DB.Where("section = ?", section).Find(&settings)
	for _, setting := range settings {
		result[setting.Key] = string(setting.Value)
	}
	return result
}

// SetSection 批量设置
func (s *SettingsService) SetSection(section string, values map[string]string) error {
	for key, value := range values {
		if err := s.Set(section, key, value); err != nil {
			return err
		}
	}
	if section == constant.SectionSite {
		cache.SetSiteCacheBatch(values)
	}
	return nil
}
