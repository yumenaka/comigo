package config

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scheduler"
)

// GlobalLibraryScanner 全局的定时扫描调度器实例
var GlobalLibraryScanner *scheduler.LibraryScanner

// scanTaskFunc 扫描任务函数，用于定时扫描（由外部包设置，避免循环导入）
var scanTaskFunc func() error

// InitLibraryScanner 初始化全局调度器
func InitLibraryScanner() {
	if GlobalLibraryScanner == nil {
		GlobalLibraryScanner = scheduler.NewLibraryScanner()
	}
}

// SetScanTaskFunc 设置扫描任务函数（由外部包调用，避免循环导入）
func SetScanTaskFunc(fn func() error) {
	scanTaskFunc = fn
}

// StartOrStopAutoRescan 根据配置启动或停止自动扫描
func StartOrStopAutoRescan() {
	InitLibraryScanner()
	interval := GetCfg().AutoRescanIntervalMinutes
	if interval > 0 {
		if scanTaskFunc == nil {
			logger.Errorf("scanTaskFunc is not set, cannot start auto rescan")
			return
		}
		if err := GlobalLibraryScanner.Start(interval, scanTaskFunc); err != nil {
			logger.Errorf(locale.GetString("log_scheduler_create_scheduler_failed"), err)
		} else {
			logger.Infof(locale.GetString("auto_rescan_started"), interval)
		}
	} else {
		// 如果间隔为 0，停止定时扫描
		if err := GlobalLibraryScanner.Stop(); err != nil {
			logger.Errorf(locale.GetString("log_scheduler_stop_task_failed"), err)
		} else {
			logger.Infof(locale.GetString("auto_rescan_stopped"))
		}
	}
}
