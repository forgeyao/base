package logger

import (
	"errors"
	"fmt"
)

// Entry is a single logger config. Unset fields inherit from Config.
type Entry struct {
	Name       string `yaml:"name"`
	Filename   string `yaml:"filename"`
	Format     string `yaml:"format,omitempty"`     // "raw" outputs only the message body.
	Level      string `yaml:"level,omitempty"`      // Overrides the shared level.
	MaxSize    int    `yaml:"maxsize,omitempty"`    // Overrides the shared MaxSize in MB.
	MaxAge     int    `yaml:"maxage,omitempty"`     // Overrides the shared MaxAge in days.
	MaxBackups int    `yaml:"maxbackups,omitempty"` // Overrides the shared MaxBackups.
}

// Config defines configuration used by the logger package.
type Config struct {
	// Shared defaults inherited by each logger entry unless overridden.
	Level      string `yaml:"level,omitempty"`      // debug, info, warn, error...
	MaxSize    int    `yaml:"maxsize,omitempty"`    // Max size for a single log file in MB.
	MaxAge     int    `yaml:"maxage,omitempty"`     // Number of days to retain log files.
	MaxBackups int    `yaml:"maxbackups,omitempty"` // Maximum number of rotated backups.

	// Multi-logger config. The first name="" entry is treated as the global Log.
	Loggers []Entry `yaml:"loggers,omitempty"`

	// Legacy single-file mode. Ignored when Loggers is not empty.
	Filename string `yaml:"filename,omitempty"`
}

// Validate normalizes and validates the logger config.
func (c *Config) Validate() error {
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

	if c.Filename == "" {
		return errors.New("conf.Log.Filename must provided")
	}
	return nil
}
