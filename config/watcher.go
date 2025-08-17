package config

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher 配置文件监听器
type Watcher struct {
	watcher   *fsnotify.Watcher
	config    *Config
	callbacks []WatchCallback
	mu        sync.RWMutex
	stopCh    chan struct{}
	running   bool
}

// NewWatcher 创建新的配置文件监听器
func NewWatcher(config *Config) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("创建文件监听器失败: %w", err)
	}

	return &Watcher{
		watcher:   watcher,
		config:    config,
		callbacks: make([]WatchCallback, 0),
		stopCh:    make(chan struct{}),
		running:   false,
	}, nil
}

// AddCallback 添加配置变化回调
func (w *Watcher) AddCallback(callback WatchCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks = append(w.callbacks, callback)
}

// Start 开始监听配置文件
func (w *Watcher) Start(configPath string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("监听器已经在运行")
	}

	// 添加配置文件到监听列表
	err := w.watcher.Add(configPath)
	if err != nil {
		return fmt.Errorf("添加文件监听失败: %w", err)
	}

	// 同时监听配置文件所在的目录（处理文件重命名等情况）
	dir := filepath.Dir(configPath)
	err = w.watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("添加目录监听失败: %w", err)
	}

	w.running = true

	// 启动监听协程
	go w.watchLoop(configPath)

	return nil
}

// Stop 停止监听
func (w *Watcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return nil
	}

	w.running = false
	close(w.stopCh)

	return w.watcher.Close()
}

// watchLoop 监听循环
func (w *Watcher) watchLoop(configPath string) {
	// 防抖动：短时间内的多次事件只处理一次
	debounceTimer := time.NewTimer(0)
	debounceTimer.Stop()
	
	var pendingReload bool

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// 只处理配置文件的写入和创建事件
			if w.shouldReload(event, configPath) {
				// 设置防抖动定时器
				debounceTimer.Reset(100 * time.Millisecond)
				pendingReload = true
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("配置文件监听错误: %v\n", err)

		case <-debounceTimer.C:
			if pendingReload {
				w.handleConfigChange(configPath)
				pendingReload = false
			}

		case <-w.stopCh:
			debounceTimer.Stop()
			return
		}
	}
}

// shouldReload 判断是否应该重新加载配置
func (w *Watcher) shouldReload(event fsnotify.Event, configPath string) bool {
	// 检查是否是目标配置文件
	if event.Name != configPath {
		// 如果是目录中的文件，检查是否是配置文件
		if filepath.Dir(event.Name) == filepath.Dir(configPath) {
			if filepath.Base(event.Name) == filepath.Base(configPath) {
				return true
			}
		}
		return false
	}

	// 只处理写入、创建和重命名事件
	return event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create ||
		event.Op&fsnotify.Rename == fsnotify.Rename
}

// handleConfigChange 处理配置变化
func (w *Watcher) handleConfigChange(configPath string) {
	// 保存旧配置的副本
	oldConfig := w.copyConfig(w.config.data)

	// 重新加载配置
	loader := NewLoader(w.config)
	err := loader.LoadFromFile(configPath)
	if err != nil {
		fmt.Printf("重新加载配置文件失败: %v\n", err)
		return
	}

	// 加载环境变量覆盖
	envManager := NewEnvManager(w.config)
	envManager.LoadEnvVars()

	// 获取新配置
	newConfig := w.copyConfig(w.config.data)

	// 调用所有回调函数
	w.mu.RLock()
	callbacks := make([]WatchCallback, len(w.callbacks))
	copy(callbacks, w.callbacks)
	w.mu.RUnlock()

	for _, callback := range callbacks {
		go func(cb WatchCallback) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("配置变化回调函数执行出错: %v\n", r)
				}
			}()
			cb(oldConfig, newConfig)
		}(callback)
	}
}

// copyConfig 深拷贝配置数据
func (w *Watcher) copyConfig(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	result := make(map[string]interface{})
	for key, value := range data {
		result[key] = w.copyValue(value)
	}
	return result
}

// copyValue 深拷贝值
func (w *Watcher) copyValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return w.copyConfig(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = w.copyValue(item)
		}
		return result
	default:
		return v
	}
}

// IsRunning 检查监听器是否正在运行
func (w *Watcher) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

// GetCallbackCount 获取回调函数数量
func (w *Watcher) GetCallbackCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.callbacks)
}

// ClearCallbacks 清除所有回调函数
func (w *Watcher) ClearCallbacks() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks = w.callbacks[:0]
}

// RemoveCallback 移除指定的回调函数（通过索引）
func (w *Watcher) RemoveCallback(index int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if index < 0 || index >= len(w.callbacks) {
		return fmt.Errorf("回调函数索引超出范围: %d", index)
	}

	// 移除指定索引的回调函数
	w.callbacks = append(w.callbacks[:index], w.callbacks[index+1:]...)
	return nil
}
