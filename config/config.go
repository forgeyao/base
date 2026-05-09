package config

import (
	"errors"
	"fmt"
)

// LogEntry 单个 logger 的配置，未填写的字段继承 Log 中的公共配置
type LogEntry struct {
	Name       string `yaml:"name"`
	Filename   string `yaml:"filename"`
	Format     string `yaml:"format,omitempty"`     // 输出格式："raw" = 只输出消息，无时间/级别/caller 前缀；默认为完整格式
	Level      string `yaml:"level,omitempty"`      // 覆盖公共 Level
	MaxSize    int    `yaml:"maxsize,omitempty"`    // 覆盖公共 MaxSize（MB）
	MaxAge     int    `yaml:"maxage,omitempty"`     // 覆盖公共 MaxAge（天）
	MaxBackups int    `yaml:"maxbackups,omitempty"` // 覆盖公共 MaxBackups
}

type Log struct {
	// 公共配置，所有 logger 共享，可被 LogEntry 中的同名字段覆盖
	Level      string `yaml:"level,omitempty"`      // 日志等级: debug, info, warn, error...
	MaxSize    int    `yaml:"maxsize,omitempty"`    // 单个日志文件大小上限（MB）
	MaxAge     int    `yaml:"maxage,omitempty"`     // 日志文件保留天数
	MaxBackups int    `yaml:"maxbackups,omitempty"` // 最大备份文件数

	// 多 logger 配置；优先使用第一个 name="" 的 logger 作为全局 logger.Log
	Loggers []LogEntry `yaml:"loggers,omitempty"`

	// 兼容旧版单文件配置，Loggers 不为空时忽略此字段
	Filename string `yaml:"filename,omitempty"`
}

type Validator interface {
	Validate() error
}

func (c *Log) Validate() error {
	// 公共默认值
	if c.MaxSize <= 0 {
		c.MaxSize = 100
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 7
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 7
	}

	if len(c.Loggers) > 0 {
		// 多 logger 模式：每个 entry 必须有 name 和 filename
		for i, entry := range c.Loggers {
			if entry.Name == "" {
				return fmt.Errorf("conf.Log.Loggers[%d].name must provided", i)
			}
			if entry.Filename == "" {
				return fmt.Errorf("conf.Log.Loggers[%d].filename must provided", i)
			}
		}
		return nil
	}

	// 兼容旧版单文件模式
	if c.Filename == "" {
		return errors.New("conf.Log.Filename must provided")
	}
	return nil
}
