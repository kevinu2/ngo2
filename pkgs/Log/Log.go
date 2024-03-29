package Log

import (
	"fmt"
	"github.com/kevinu2/ngo2/pkgs/Default"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var l *Log

type Log struct {
	SugaredLogger *zap.SugaredLogger
	Config        *Config
}

func init() {
	l = New()
	DefaultConfig()
}

func New() *Log {
	return new(Log)
}

func AddConfig(logLevel, logPath, logFile, logOutput string) {
	l.AddConfig(logLevel, logPath, logFile, logOutput)
}
func (l *Log) AddConfig(logLevel, logPath, logFile, logOutput string) {
	l.Config = &Config{
		LogLevel:  logLevel,
		LogPath:   logPath,
		LogFile:   logFile,
		LogOutput: logOutput,
	}
}

func DefaultConfig() {
	l.Config = &Config{
		LogLevel:  "Info",
		LogPath:   "/tmp",
		LogFile:   "app.log",
		LogOutput: "std",
	}
}

func Logger() *zap.SugaredLogger {
	return l.GetLogger()
}
func (l *Log) GetLogger() *zap.SugaredLogger {
	if l.SugaredLogger == nil {
		fmt.Print("Log: initLogger()! \n")
		l.initLoggerLevel()
	}
	return l.SugaredLogger
}

func (l *Log) initLogger(level zapcore.Level) {
	var core zapcore.Core

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(Default.TimeUtcFormat))
		},
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	switch l.Config.LogOutput {
	case LogOutputStd.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	case LogOutputFile.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(getWriter(fmt.Sprintf("%s/%s", l.Config.LogPath, l.Config.LogFile))), level),
		)
		//logger := zap.NewDevelopment()
	case LogOutputBoth.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
			zapcore.NewCore(encoder, zapcore.AddSync(getWriter(fmt.Sprintf("%s/%s", l.Config.LogPath, l.Config.LogFile))), level),
		)
	default:
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	}
	l.SugaredLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	return
}

func (l *Log) initLoggerLevel() {
	switch l.Config.LogLevel {
	case LogDebug.Level():
		l.initLogger(zapcore.DebugLevel)
		return
	case LogInfo.Level():
		l.initLogger(zapcore.InfoLevel)
		return
	case LogWarn.Level():
		l.initLogger(zapcore.WarnLevel)
		return
	case LogError.Level():
		l.initLogger(zapcore.ErrorLevel)
		return
	case LogDPanic.Level():
		l.initLogger(zapcore.DPanicLevel)
		return
	case LogPanic.Level():
		l.initLogger(zapcore.PanicLevel)
		return
	case LogFatal.Level():
		l.initLogger(zapcore.FatalLevel)
		return
	default:
		l.initLogger(zapcore.InfoLevel)
		return
	}
}

func getWriter(filename string) io.Writer {
	hook, err := rotates.New(
		filename+".%Y%m%d%H",
		rotates.WithLinkName(filename),
		rotates.WithMaxAge(time.Hour*24*7),
		rotates.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
