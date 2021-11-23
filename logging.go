package logging
import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
)

var ZapLogger *zap.Logger

func InitLogger(log LogConf) {
	log.initDefault()
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level", // 输入日志级别的key名
		MessageKey:     "msg", // 输出信息的key名
		NameKey:        "logger",
		CallerKey:      "file",//"caller"
		StacktraceKey:  "stacktrace",//"trace"
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		//EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短门路编码器
		EncodeName:     zapcore.FullNameEncoder,
	})
	// 设置日志级别
	//atomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(zap.DebugLevel)
	// 设置日志级别 自定义日志等级的interface
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= log.LogLevel
	})
	var writes []zapcore.Core
	writes = append(writes, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(log)),level))
	// 如果是开发环境，同时在管制台上也输入
	if log.Debug != "prod" {
		writes = append(writes, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}
	core := zapcore.NewTee(writes...)
	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	development := zap.Development()
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel),development)
	defer ZapLogger.Sync()
}
func getWriter(log LogConf) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		// 没有使用go风格反人类的format格式
		log.Path + log.Name + "-" + log.FileName + log.Suffix,
		//log.Path + log.Name + ".%Y-%m-%d %H:%M:%S.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(log.Path + log.Name + log.Suffix),
		rotatelogs.WithMaxAge(time.Hour * time.Duration(log.MaxSaveDay)),
		rotatelogs.WithRotationTime(time.Hour * time.Duration(log.CutHour)),
	)
	if err != nil {
		panic(err)
	}
	return hook
}