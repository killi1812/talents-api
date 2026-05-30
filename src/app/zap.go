package app

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: change setups so taht zap logs and fmt prints to Stdout

func prodLoggerSetup() error {
	consoleLogLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.InfoLevel
	})

	// log output
	consoleLogFile := zapcore.Lock(os.Stdout)

	// log configuration
	consoleLogConfig := zap.NewProductionEncoderConfig()
	consoleLogConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleLogConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleLogEncoder := zapcore.NewConsoleEncoder(consoleLogConfig)

	// setup zap
	// duplicate log entries into multiple cores
	core := zapcore.NewTee(
		zapcore.NewCore(consoleLogEncoder, consoleLogFile, consoleLogLevel),
	)

	// create logger from core
	// options = annotate message with the filename, line number, and function name
	logger := zap.New(core, zap.AddCaller())
	// No defer Sync here since it closes os.Stdout

	// replace global logger
	_ = zap.ReplaceGlobals(logger)

	return nil
}

func devLoggerSetup() error {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		return err
	}
	_ = zap.ReplaceGlobals(logger)

	return nil
}

func VerboseLoggerSetup() error {
	consoleLogLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.DebugLevel
	})

	// log output
	consoleLogFile := zapcore.Lock(os.Stdout)

	// log configuration
	consoleLogConfig := zap.NewProductionEncoderConfig()
	consoleLogConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleLogConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleLogEncoder := zapcore.NewConsoleEncoder(consoleLogConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleLogEncoder, consoleLogFile, consoleLogLevel),
	)

	// create logger from core
	logger := zap.New(core, zap.AddCaller())
	// defer logger.Sync()

	zap.S().Debugln("New verbose logger created successfully")
	// replace global logger
	_ = zap.ReplaceGlobals(logger)

	zap.S().Debugln("New verbose logger replace successfully")
	return nil
}
