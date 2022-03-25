package mr_smart

import (
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap"
)

var serverLogger = util.NewLogger("./server.log", 10, 10, 7, zap.WarnLevel, false)
