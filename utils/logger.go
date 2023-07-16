package utils

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *zap.SugaredLogger
	loggerOnce sync.Once
)

func init() {
    logger= InitLogger()
}

func InitLogger() *zap.SugaredLogger {
    loggerOnce.Do(func() {
        var err error
        var customLogger *zap.Logger

        cfg := zap.NewProductionConfig()
        cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
    
        if customLogger, err = cfg.Build(); err != nil {
            panic(fmt.Errorf("logger init faild"))
        }
        logger = customLogger.Sugar()
    })
	return logger
}
