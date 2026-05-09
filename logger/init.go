package logger

import (
	"time"

	conf "util/config"

	"github.com/DeRuina/timberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局默认 logger，指向 Loggers 中 name="" 的实例。
var Log *zap.SugaredLogger

var loggers = map[string]*zap.SugaredLogger{}

// InitLogger 兼容单文件配置。
func InitLogger(l conf.Log) {
	writer := newWriter(l.Filename, l.MaxSize, l.MaxAge, l.MaxBackups)
	logger := buildLogger(writer, l.Level, "")
	Log = logger
	loggers[""] = logger
	Log.Info("Log init success")
	//nolint:errcheck
	defer Log.Sync()
}

// InitLoggers 多 logger 初始化，未填写的字段继承 l 中的公共配置。
// 优先使用 name="" 的 logger 作为全局 Log，否则使用第一个 logger。
func InitLoggers(l conf.Log) {
	loggers = map[string]*zap.SugaredLogger{}

	for _, entry := range l.Loggers {
		maxSize := entry.MaxSize
		if maxSize <= 0 {
			maxSize = l.MaxSize
		}
		maxAge := entry.MaxAge
		if maxAge <= 0 {
			maxAge = l.MaxAge
		}
		maxBackups := entry.MaxBackups
		if maxBackups <= 0 {
			maxBackups = l.MaxBackups
		}
		level := entry.Level
		if level == "" {
			level = l.Level
		}

		writer := newWriter(entry.Filename, maxSize, maxAge, maxBackups)
		logger := buildLogger(writer, level, entry.Format)
		loggers[entry.Name] = logger
	}

	if lg, ok := loggers[""]; ok {
		Log = lg
	} else if len(l.Loggers) > 0 {
		Log = loggers[l.Loggers[0].Name]
	}

	if Log != nil {
		Log.Info("Log init success, loggers:", loggerNames(l.Loggers))
		//nolint:errcheck
		defer Log.Sync()
	}
}

// GetLogger 按名称获取 logger，name 不存在时返回全局 Log
func GetLogger(name string) *zap.SugaredLogger {
	if lg, ok := loggers[name]; ok {
		return lg
	}
	return Log
}

func namedLogger(name string) *zap.SugaredLogger {
	if lg := GetLogger(name); lg != nil {
		return lg
	}
	return zap.NewNop().Sugar()
}

// Info 记录信息日志。
func Info(format string, v ...interface{}) {
	namedLogger("").Infof(format, v...)
}

// Error 记录错误日志。
func Error(format string, v ...interface{}) {
	defaultLogger := namedLogger("")
	errorLogger := namedLogger("error")

	defaultLogger.Errorf(format, v...)
	if errorLogger != defaultLogger {
		errorLogger.Errorf(format, v...)
	}
}

func loggerNames(entries []conf.LogEntry) []string {
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name)
	}
	return names
}

func buildLogger(writer zapcore.WriteSyncer, levelStr string, format string) *zap.SugaredLogger {
	level := zapcore.InfoLevel
	if levelStr != "" {
		var err error
		level, err = zapcore.ParseLevel(levelStr)
		if err != nil {
			level = zapcore.InfoLevel
		}
	}
	if format == "raw" {
		// raw 模式：只输出消息本身，不加时间/级别/caller 前缀
		enc := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		})
		core := zapcore.NewCore(enc, writer, level)
		return zap.New(core).Sugar()
	}
	core := zapcore.NewCore(getEncoder(), writer, level)
	return zap.New(core, zap.AddCaller()).Sugar()
}

func newWriter(filename string, maxSize, maxAge, maxBackups int) zapcore.WriteSyncer {
	lumber := &timberjack.Logger{
		Filename:         filename,
		MaxSize:          maxSize,
		MaxAge:           maxAge,
		MaxBackups:       maxBackups,
		LocalTime:        true,
		Compression:      "gzip",
		RotateAt:         []string{"00:00"}, // 每天零点轮转
		BackupTimeFormat: "2006-01-02_15-04-05",
		FileMode:         0644,
	}
	return zapcore.AddSync(lumber)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(NewEncoderConfig())
}

// https://pkg.go.dev/go.uber.org/zap@v1.21.0/zapcore#EncoderConfig
func NewEncoderConfig() zapcore.EncoderConfig {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// https://pkg.go.dev/go.uber.org/zap@v1.21.0#Config
func NewConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    NewEncoderConfig(),
		OutputPaths:      []string{"stdout", "server.log"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
