package utility

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var srvLogger = NewLogger("./server.log", 10, 10, 7, zap.DebugLevel, false)

// SetLogLevel must greater than debug level
func SetLogLevel(lv zapcore.Level) {
	srvLogger = srvLogger.WithOptions(zap.IncreaseLevel(lv))
}

func DefaultLogger() *zap.Logger {
	return srvLogger
}

func NewLogger(logFile string, maxFileMB, maxBackups, maxDays int, level zapcore.Level, dev bool) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logFile,    // 日志文件路径
		MaxSize:    maxFileMB,  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxDays,    // 文件最多保存多少天
		Compress:   true,       // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	//atomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(zap.InfoLevel)
	//
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zap.NewAtomicLevelAt(level),                                                     // 日志级别
	)
	var opts = []zap.Option{zap.AddStacktrace(zap.ErrorLevel)}
	//
	if dev {
		// 开启开发模式，堆栈跟踪
		opts = append(opts, zap.AddCaller())
		//caller := zap.AddCaller()
		// 开启文件及行号
		opts = append(opts, zap.Development())
		//development := zap.Development()
		// 设置初始化字段
		//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	}
	// 构造日志
	return zap.New(core, opts...)
	//
	//logger.Info("log 初始化成功")
	//logger.Info("无法获取网址",
	//	zap.String("url", "http://www.baidu.com"),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second))
}
