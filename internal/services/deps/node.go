package deps

import (
	"encoding/json"

	"github.com/hicongcn/xuanwu-panel/internal/logger"
	"github.com/hicongcn/xuanwu-panel/internal/models"
)

type NodeManager struct {
	BaseManager
}

func NewNodeManager(language string) *NodeManager {
	return &NodeManager{
		BaseManager: BaseManager{
			Language:     language,
			InstallCmd:   []string{"npm", "install", "-g"},
			UninstallCmd: []string{"npm", "uninstall", "-g"},
			ListCmd:      []string{"npm", "list", "-g", "--depth=0", "--json"},
			VerifyCmd:    []string{"node", "-v"},
			Separator:    "@",
		},
	}
}

type npmListOutput struct {
	Dependencies map[string]struct {
		Version string `json:"version"`
	} `json:"dependencies"`
}

func (m *NodeManager) GetInstalledPackages(language, langVersion string) ([]models.Dependency, error) {
	output, err := m.runMiseCommand(langVersion, m.ListCmd)
	if err != nil {
		logger.Warnf("GetInstalledPackages for %s failed: %v", language, err)
	}

	var npmOut npmListOutput
	if err := json.Unmarshal(output, &npmOut); err != nil {
		logger.Warnf("GetInstalledPackages for %s: failed to parse npm output: %v", language, err)
		return nil, err
	}

	var packages []models.Dependency
	for name, info := range npmOut.Dependencies {
		packages = append(packages, models.Dependency{
			Name:     name,
			Version:  info.Version,
			Language: language,
		})
	}
	return packages, nil
}
