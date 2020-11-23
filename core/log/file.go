package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"x-server/core/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const UniqueId = "log_id"
const DefaultRequest = "-"
const DefaultRorate = 24 * time.Hour

var (
	logger *zap.SugaredLogger
)

type Log struct {
	Local LogLocal
}

type LogLocal struct {
	Enable bool             `yaml:"enable"`
	Config []LogLocalConfig `yaml:"config"`
}

type LogLocalConfig struct {
	Path       string `yaml:"path"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max-size"`
	MaxBackups int    `yaml:"max-backups"`
	MaxAge     int    `yaml:"max-age"`
	LevelMin   int8   `yaml:"level-min"`
	LevelMax   int8   `yaml:"level-max"`
}

func New() error {
	var logConfig = &Log{}
	config.Viper.UnmarshalKey("log.log", logConfig)

	var cores []zapcore.Core
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	config.Viper.Get("log.local.config")
	if true {
		for _, c := range logConfig.Local.Config {
			_, err := os.Stat(c.Path)
			if os.IsNotExist(err) {
				err := os.Mkdir(c.Path, os.ModePerm)
				if err != nil {
					return err
				}
			}
			ll := &lumberjack.Logger{
				Filename:   c.Path + "/" + c.File,
				MaxSize:    c.MaxSize,
				MaxBackups: c.MaxBackups,
				MaxAge:     c.MaxAge, // days
				LocalTime:  true,
			}
			go rorate(ll)
			levelMax := zapcore.Level(c.LevelMax)
			levelMin := zapcore.Level(c.LevelMin)
			enableFunc := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev <= levelMax && lev >= levelMin
			})
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(ll), enableFunc))
		}
	}

	core := zapcore.NewTee(cores...)
	logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()

	return nil
}

func rorate(l *lumberjack.Logger) {
	now := time.Now()
	_, offset := now.Zone()
	renewDur := now.Round(DefaultRorate)
	if renewDur.Before(now) {
		renewDur = renewDur.Add(DefaultRorate)
	}

	renewTick := time.NewTimer(renewDur.Sub(now.Add(time.Duration(offset) * time.Second)))
	for {
		select {
		case <-renewTick.C:
			l.Rotate()
			renewTick.Reset(DefaultRorate)
		case <-context.Background().Done():
			Info(nil, "log rorate done")
		}
	}
}

func Debugf(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.DebugLevel, msg, output...)
	if logger != nil {
		logger.Debugw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Infof(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.InfoLevel, msg, output...)
	if logger != nil {
		logger.Infow(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Warnf(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.WarnLevel, msg, output...)
	if logger != nil {
		logger.Warnw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Errorf(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.ErrorLevel, msg, output...)
	if logger != nil {
		logger.Errorw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Fatalf(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.FatalLevel, msg, output...)
	if logger != nil {
		logger.Fatalw(msg, stack)
	} else {
		log.Fatal(msg)
	}
}

func Panicf(ctx context.Context, msg string, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.PanicLevel, msg, output...)
	if logger != nil {
		logger.Panicw(msg, stack)
	} else {
		log.Panic(msg)
	}
}

func Debug(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.DebugLevel, "", output...)
	if logger != nil {
		logger.Debugw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Info(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.InfoLevel, "", output...)
	if logger != nil {
		logger.Infow(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Warn(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.WarnLevel, "", output...)
	if logger != nil {
		logger.Warnw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Error(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.ErrorLevel, "", output...)
	if logger != nil {
		logger.Errorw(msg, stack)
	} else {
		log.Print(msg)
	}
}

func Fatal(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.FatalLevel, "", output...)
	if logger != nil {
		logger.Fatalw(msg, stack)
	} else {
		log.Fatal(msg)
	}
}

func Panic(ctx context.Context, output ...interface{}) {
	msg, stack := logMsg(ctx, zap.PanicLevel, "", output...)
	if logger != nil {
		logger.Panicw(msg, stack)
	} else {
		log.Panic(msg)
	}
}

func logMsg(ctx context.Context, lvl zapcore.Level, msg string, output ...interface{}) (string, zap.Field) {
	if msg == "" && len(output) > 0 {
		msg = fmt.Sprint(output...)
	} else if msg != "" && len(output) > 0 {
		msg = fmt.Sprintf(msg, output...)
	}
	logId := DefaultRequest
	if ctx != nil {
		if id, ok := ctx.Value(UniqueId).(uint64); ok {
			logId = strconv.FormatUint(id, 10)
		}
	}
	msg = logId + "\t" + msg
	msg = strings.Replace(msg, "\n", "", -1)
	if zap.WarnLevel.Enabled(lvl) {
		return msg, zap.Stack("stack")
	}
	return msg + "\t-", zap.Skip()
}
