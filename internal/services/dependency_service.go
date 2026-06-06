package services

import (
	"errors"

	"github.com/hicongcn/xuanwu-panel/internal/database"
	"github.com/hicongcn/xuanwu-panel/internal/logger"
	"github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/services/deps"
	"github.com/hicongcn/xuanwu-panel/internal/utils"
)

type DependencyService struct{}

func NewDependencyService() *DependencyService {
	return &DependencyService{}
}

// List 获取依赖列表，同时同步实际安装的包到数据库
func (s *DependencyService) List(language, langVersion string) ([]models.Dependency, error) {
	// 先同步实际安装的包到数据库（补充版本号和未入库的包）
	s.syncInstalledPackages(language, langVersion)

	var results []models.Dependency
	query := database.DB
	if language != "" {
		query = query.Where("language = ?", language)
	}
	if langVersion != "" {
		query = query.Where("lang_version = ?", langVersion)
	}
	err := query.Order("id desc").Find(&results).Error
	return results, err
}

// syncInstalledPackages 将实际安装的包同步到数据库，补全版本号
func (s *DependencyService) syncInstalledPackages(language, langVersion string) {
	if language == "" {
		return
	}
	installed, err := s.GetInstalledPackages(language, langVersion)
	if err != nil {
		logger.Warnf("syncInstalledPackages failed for %s: %v", language, err)
		return
	}
	if len(installed) == 0 {
		return
	}

	// 构建已安装包的版本映射
	installedMap := make(map[string]string)
	for _, pkg := range installed {
		installedMap[pkg.Name] = pkg.Version
	}

	// 查询数据库中该语言下的所有包
	var dbDeps []models.Dependency
	query := database.DB.Where("language = ?", language)
	if langVersion != "" {
		query = query.Where("lang_version = ?", langVersion)
	}
	if err := query.Find(&dbDeps).Error; err != nil {
		return
	}

	// 更新数据库中版本为空的记录
	for i := range dbDeps {
		if dbDeps[i].Version == "" {
			if v, ok := installedMap[dbDeps[i].Name]; ok && v != "" {
				_ = database.DB.Model(&dbDeps[i]).Update("version", v).Error
			}
		}
	}

	// 将数据库中不存在的已安装包入库
	dbNames := make(map[string]bool)
	for _, d := range dbDeps {
		dbNames[d.Name] = true
	}
	for name, version := range installedMap {
		if !dbNames[name] && version != "" {
			dep := models.Dependency{
				Name:        name,
				Version:     version,
				Language:    language,
				LangVersion: langVersion,
			}
			dep.ID = utils.GenerateID()
			_ = database.DB.Create(&dep).Error
		}
	}
}

// Create 创建依赖记录
func (s *DependencyService) Create(dep *models.Dependency) error {
	// 检查是否已存在（名称、版本、语言及版本必须完全匹配）
	var existing models.Dependency
	res := database.DB.Where("name = ? AND version = ? AND language = ? AND lang_version = ?", dep.Name, dep.Version, dep.Language, dep.LangVersion).Limit(1).Find(&existing)
	if res.Error == nil && res.RowsAffected > 0 {
		// 如果已存在，更新 ID 并执行更新
		dep.ID = existing.ID
		return database.DB.Model(&existing).Updates(dep).Error
	}

	// 不存在则新建
	if dep.ID == "" {
		dep.ID = utils.GenerateID()
	}
	return database.DB.Create(dep).Error
}

// Delete 删除依赖记录
func (s *DependencyService) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Dependency{}).Error
}

// Install 安装依赖
func (s *DependencyService) Install(dep *models.Dependency) error {
	m := deps.GetManager(dep.Language)
	if m == nil {
		return errors.New("不支持的依赖类型: " + dep.Language)
	}
	return m.Install(dep)
}

// Uninstall 卸载依赖
func (s *DependencyService) Uninstall(dep *models.Dependency) error {
	m := deps.GetManager(dep.Language)
	if m == nil {
		return errors.New("不支持的依赖类型: " + dep.Language)
	}
	return m.Uninstall(dep)
}

// GetInstalledPackages 获取已安装的包列表
func (s *DependencyService) GetInstalledPackages(language, langVersion string) ([]models.Dependency, error) {
	m := deps.GetManager(language)
	if m == nil {
		return nil, errors.New("不支持的依赖类型: " + language)
	}
	return m.GetInstalledPackages(language, langVersion)
}

// GetInstallCommand 获取安装命令
func (s *DependencyService) GetInstallCommand(dep *models.Dependency) (string, error) {
	m := deps.GetManager(dep.Language)
	if m == nil {
		return "", errors.New("不支持的依赖类型: " + dep.Language)
	}
	return m.GetInstallCommand(dep)
}

// GetReinstallAllCommand 获取全部重装命令
func (s *DependencyService) GetReinstallAllCommand(language, langVersion string) (string, error) {
	m := deps.GetManager(language)
	if m == nil {
		return "", errors.New("不支持的依赖类型: " + language)
	}

	deps_list, err := s.List(language, langVersion)
	if err != nil {
		return "", err
	}

	return m.GetReinstallAllCommand(deps_list)
}

// GetBatchInstallCommand 获取批量安装命令
func (s *DependencyService) GetBatchInstallCommand(depsList []models.Dependency) (string, error) {
	if len(depsList) == 0 {
		return "", errors.New("依赖包列表不能为空")
	}

	firstDep := depsList[0]
	m := deps.GetManager(firstDep.Language)
	if m == nil {
		return "", errors.New("不支持的依赖类型: " + firstDep.Language)
	}

	return m.GetBatchInstallCommand(depsList)
}

// ImportDependencies 批量导入依赖并自动入库去重
func (s *DependencyService) ImportDependencies(depsList []models.Dependency) ([]models.Dependency, error) {
	var imported []models.Dependency
	for i := range depsList {
		dep := &depsList[i]
		if err := s.Create(dep); err == nil {
			imported = append(imported, *dep)
		}
	}
	return imported, nil
}
