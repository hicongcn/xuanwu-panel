package controllers

import (
	"fmt"

	"github.com/hicongcn/xuanwu-panel/internal/constant"
	"github.com/hicongcn/xuanwu-panel/internal/database"
	"github.com/hicongcn/xuanwu-panel/internal/models"
	"github.com/hicongcn/xuanwu-panel/internal/services/tasks"
	"github.com/hicongcn/xuanwu-panel/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type LogWSController struct{}

func NewLogWSController() *LogWSController {
	return &LogWSController{}
}

func (lc *LogWSController) StreamLog(c *gin.Context) {
	logIDStr := c.Query("log_id")
	if logIDStr == "" {
		return
	}

	logID := logIDStr

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// DoS 保护与心跳设置
	conn.SetReadLimit(constant.MaxMessageSize) // 使用与终端一致的限制
	conn.SetReadDeadline(time.Now().Add(constant.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(constant.PongWait))
		return nil
	})

	// 启动一个读取循环，用于处理 pong 和检测断开
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()

	// 1. 检查数据库中是否已结束
	var taskLog models.TaskLog
	res := database.DB.Where("id = ?", logID).Limit(1).Find(&taskLog)
	if res.Error == nil && res.RowsAffected > 0 {
		if taskLog.Status != "running" {
			// 已结束，读取库内日志
			content, err := utils.DecompressFromBase64(string(taskLog.Output))
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("解压日志失败: "+err.Error()))
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte(content))
			return
		}
	}

	// 2. 未结束或未找到记录，尝试从 TinyLogManager 获取
	tl := tasks.GetActiveLog(logID)
	if tl == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("未找到正在运行的任务日志"))
		return
	}

	// 发送系统提示
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[System] 连接成功，正在监听日志... (LogID: %s)\n", logID)))

	// 发送最后 100 行
	lastLines, err := tl.ReadLastLines(100)
	if err == nil && len(lastLines) > 0 {
		conn.WriteMessage(websocket.TextMessage, lastLines)
	}

	// 订阅实时更新
	sub := tl.Subscribe()
	defer tl.Unsubscribe(sub)

	ticker := time.NewTicker(constant.PingPeriod)
	defer ticker.Stop()

	// 推送更新
	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case data, ok := <-sub:
			if !ok {
				// 任务结束，尝试刷新最后一次库内完整内容
				var finalLog models.TaskLog
				res := database.DB.Where("id = ?", logID).Limit(1).Find(&finalLog)
				if res.Error == nil && res.RowsAffected > 0 {
					content, _ := utils.DecompressFromBase64(string(finalLog.Output))
					if content != "" {
						conn.WriteMessage(websocket.TextMessage, []byte("\n--- 任务已结束 ---\n"))
						// 这里可以选择性再推一次完整版，或直接退出
					}
				}
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}
