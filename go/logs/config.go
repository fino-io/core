package logs

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	InitFields map[string]any `json:"initFields"`
	Level      string         `json:"level" default:"debug"`    // debug,info,warn,error,fatal
	Encode     string         `json:"encode" default:"console"` // console,json
	Output     string         `json:"output" default:"console"` // console,file
	File       FileConfig     `json:"file"`
}

type FileConfig struct {
	Path       string `json:"path" default:"./logs/app.log"`
	Encode     string `json:"encode" default:"json"`
	MaxSize    int    `json:"maxSize" default:"100"`
	MaxBackups int    `json:"maxBackups" default:"10"`
	MaxAge     int    `json:"maxAge" default:"30"`
	Compress   bool   `json:"compress"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Level:  "debug",
		Encode: "console",
		Output: "console",
		File: FileConfig{
			Path:       "./logs/app.log",
			Encode:     "json",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
		},
	}
}

// Validate accepts omitted fields as defaults, but rejects invalid settings.
func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	if c.Level != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(c.Level)); err != nil {
			return fmt.Errorf("invalid log level %q: %w", c.Level, err)
		}
	}
	for name, value := range map[string]string{"encode": c.Encode, "file.encode": c.File.Encode} {
		if value != "" && !strings.EqualFold(value, "json") && !strings.EqualFold(value, "console") {
			return fmt.Errorf("invalid log %s: %q", name, value)
		}
	}
	if c.Output != "" && !strings.EqualFold(c.Output, "console") && !strings.EqualFold(c.Output, "file") {
		return fmt.Errorf("invalid log output: %q", c.Output)
	}
	if c.File.MaxSize < 0 || c.File.MaxAge < 0 || c.File.MaxBackups < 0 {
		return fmt.Errorf("log file limits must be non-negative")
	}
	return nil
}

func (c *Config) withDefaults() *Config {
	defaults := NewDefaultConfig()
	if c == nil {
		return defaults
	}
	resolved := *c
	if resolved.Level == "" {
		resolved.Level = defaults.Level
	}
	if resolved.Encode == "" {
		resolved.Encode = defaults.Encode
	}
	if resolved.Output == "" {
		resolved.Output = defaults.Output
	}
	if resolved.File.Path == "" {
		resolved.File.Path = defaults.File.Path
	}
	if resolved.File.Encode == "" {
		resolved.File.Encode = defaults.File.Encode
	}
	return &resolved
}
