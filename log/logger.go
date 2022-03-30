package log

import (
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var srvLogger = util.NewLogger("./server.log", 10, 10, 7, zap.DebugLevel, false)

// SetLogLevel must greater than debug level
func SetLogLevel(lv zapcore.Level) {
	srvLogger = srvLogger.WithOptions(zap.IncreaseLevel(lv))
}

func GetLogger() *zap.Logger {
	return srvLogger
}
