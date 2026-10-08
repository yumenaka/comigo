package scheduler

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// LibraryScanner 定时扫描书库的调度器
type LibraryScanner struct {
	scheduler       gocron.Scheduler
	mu              sync.Mutex
	running         atomic.Bool
	intervalMinutes int
}

// NewLibraryScanner 创建新的 LibraryScanner 实例
func NewLibraryScanner() *LibraryScanner {
	return &LibraryScanner{}
}

// wrapTask 用原子状态标记扫描，避免停止调度器时持锁等待任务造成死锁。
func (ls *LibraryScanner) wrapTask(taskFunc func() error) func() {
	return func() {
		if ls.running.Swap(true) {
			logger.Infof(locale.GetString("log_scheduler_task_still_running_skip"))
			return
		}
		defer ls.running.Store(false)

		if err := taskFunc(); err != nil {
			logger.Errorf(locale.GetString("log_scheduler_task_execution_failed"), err)
		} else {
			logger.Infof(locale.GetString("log_scheduler_task_execution_completed"))
		}
	}
}

// Start 启动定时扫描任务
// intervalMinutes: 执行间隔（分钟），0 表示不定时扫描
// taskFunc: 要执行的任务函数
func (ls *LibraryScanner) Start(intervalMinutes int, taskFunc func() error) error {
	if intervalMinutes == 0 {
		return ls.Stop()
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// 如果已有任务在运行，先停止旧任务
	if ls.scheduler != nil {
		if err := ls.scheduler.Shutdown(); err != nil {
			logger.Errorf(locale.GetString("log_scheduler_stop_old_task_failed"), err)
		}
		ls.scheduler = nil
	}

	// 创建新的 scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Errorf(locale.GetString("log_scheduler_create_scheduler_failed"), err)
		return err
	}

	// 创建定时任务
	// 使用 WithSingletonMode 确保任务不会并发执行
	// LimitModeReschedule 表示如果上一个任务还在运行，跳过本次执行
	wrappedTask := ls.wrapTask(taskFunc)
	_, err = s.NewJob(
		gocron.DurationJob(time.Duration(intervalMinutes)*time.Minute),
		gocron.NewTask(wrappedTask),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.WithName("library-scanner"),
	)
	if err != nil {
		logger.Errorf(locale.GetString("log_scheduler_create_task_failed"), err)
		_ = s.Shutdown()
		return err
	}

	// 启动 scheduler
	s.Start()

	ls.scheduler = s
	ls.intervalMinutes = intervalMinutes

	logger.Infof(locale.GetString("log_scheduler_task_started"), intervalMinutes)
	return nil
}

// Stop 停止定时扫描任务
func (ls *LibraryScanner) Stop() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.scheduler == nil {
		return nil
	}

	// 等待当前运行的任务完成
	// 注意：gocron 的 Shutdown 会等待正在运行的任务完成
	err := ls.scheduler.Shutdown()
	if err != nil {
		logger.Errorf(locale.GetString("log_scheduler_stop_task_failed"), err)
		return err
	}

	ls.scheduler = nil
	ls.intervalMinutes = 0

	logger.Infof(locale.GetString("log_scheduler_task_stopped"))
	return nil
}

// IsRunning 检查是否有任务正在运行
func (ls *LibraryScanner) IsRunning() bool {
	return ls.running.Load()
}

// GetInterval 获取当前的扫描间隔（分钟）
func (ls *LibraryScanner) GetInterval() int {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return ls.intervalMinutes
}
