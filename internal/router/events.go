package router

import (
	// "fmt"
	"time"

	// "github.com/hicongcn/xuanwu-panel/internal/constant"
	"github.com/hicongcn/xuanwu-panel/internal/eventbus"
	// "github.com/hicongcn/xuanwu-panel/internal/logger"
	// "github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/services"
)

func setupEventHandlers(subscribers ...eventbus.Subscriber) {
	bus := eventbus.DefaultBus

	// 遍历并统一初始化所有订阅者的事件链路
	for _, s := range subscribers {
		s.SubscribeEvents(bus)
	}
}

func startAppLogCleanup(appLogSvc *services.AppLogService) {
	// 初始化时执行一次清理
	appLogSvc.CleanUp()

	// 定期清理（每隔1小时执行一次巡检）
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		appLogSvc.CleanUp()
	}
}
