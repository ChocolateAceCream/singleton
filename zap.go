package singleton

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapOptions struct {
	LogLevel          int    // -1-debug, 0-info, 1-warn, 2-error, 3-DPanic, 4-Panic, 5-Fatal
	Development       bool   // development mode
	DisableStacktrace bool   `default:"true"`    // DisableStacktrace completely disables automatic stacktrace capturing. By default, stacktraces are captured for WarnLevel and above logs in development and ErrorLevel and above in production.
	EncodingFormat    string `default:"console"` // json or console
	Prefix            string
	EncodeLevel       string
	ServiceName       string
	OutputPath        string
}

func WithZapLogger(options ZapOptions) Option {
	return func(s *Singleton) (err error) {
		encoderConfig := GetEncoderConfig(options)
		level := zap.NewAtomicLevelAt(zapcore.Level(options.LogLevel)) // log level
		development := options.Development
		config := zap.Config{
			Level:             level,
			Development:       development,
			DisableStacktrace: options.DisableStacktrace,
			Encoding:          options.EncodingFormat,
			EncoderConfig:     encoderConfig,
			InitialFields:     map[string]interface{}{"serviceName": options.ServiceName},
			OutputPaths:       []string{"stdout", options.OutputPath},
			ErrorOutputPaths:  []string{"stderr"},
		}
		logger, err := config.Build()
		if err != nil {
			return
		}
		s.Logger = logger
		return
	}
}

func GetEncoderConfig(options ZapOptions) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     GetCustomTimeEncoder(options.Prefix), // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   GetCustomCallerEncoder, // caller trace encoder
	}
	GetEncodeLevel(&config, options.EncodeLevel)
	return
}

func GetCustomTimeEncoder(prefix string) zapcore.TimeEncoder {
	return func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(prefix + " 2006/01/02 - 15:04:05.000 "))
	}
}

func GetCustomCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

// Set logger encode level based on config
func GetEncodeLevel(c *zapcore.EncoderConfig, level string) {
	switch level {
	case "LowercaseLevelEncoder":
		c.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		c.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		c.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		c.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		c.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
}
