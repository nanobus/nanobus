package main

import (
	"context"

	"github.com/go-logr/zapr"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nanobus/nanobus/example/welcome/pkg/welcome"
)

func main() {
	ctx := context.Background()
	// Initialize logger
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	log := zapr.NewLogger(zapLog)

	app := welcome.NewApp(ctx)
	outbound := app.NewOutbound()
	service := welcome.NewService(log, outbound)

	app.RegisterInbound(service)

	if err := app.Start(); err != nil {
		log.Error(err, "Exit with error")
	}
}
