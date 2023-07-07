package core

import (
    "fmt"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var logger = newLogger()

func newLogger() *zap.SugaredLogger {
    var logger *zap.Logger
    var err error

    cfg := zap.NewProductionConfig()
    cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

    if logger, err = cfg.Build(); err != nil {
        panic(fmt.Errorf("logger init faild"))
    }

    return logger.Sugar()
}