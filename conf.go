package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConf struct {
	Path 		string		`json:"path" yaml:"path"`	// 日志存储路径 绝对路径
	Name 		string		`json:"name" yaml:"name"`	// 日志文件名
	MaxSaveDay	int64		`json:"max_save_day" yaml:"max_save_day"` // 存储最大天数 默认14天
	CutHour		int64		`json:"cut_hour" yaml:"cut_hour"`		//切割小时 //	 默认24小时 按照每天切割
	FileName	string		`json:"file_name" yaml:"file_name"`
	Suffix		string		`json:"suffix" yaml:"suffix"` // 日志后缀
	Debug		string		`json:"debug" yaml:"debug"`	// debug or prod prod时将不会控制台输出
	Level 		string		`json:"level" yaml:"level"` // 日志级别  debug->info->warn->error
	LogLevel 	zapcore.Level

}

// 初始化
func (log *LogConf)initDefault() {
	if log.Path == "" {
		log.Path = "./log/"
	}
	if log.Name == "" {
		log.Name = "app"
	}
	if log.Suffix == "" {
		log.Suffix = ".log"
	}
	if log.Debug == "" {
		log.Debug = "prod"
	}
	if log.CutHour == 0 {
		log.CutHour = 24
	}
	if log.MaxSaveDay == 0 {
		log.MaxSaveDay = 14
	}
	if log.FileName == "" {
		log.FileName = "%Y-%m-%d" //%Y-%m-%d-%H-%M
	}

	// 默认打印info级别 debug->info->warn->error
	switch log.Level {
	case "debug":
		log.LogLevel = zap.DebugLevel
	case "info":
		log.LogLevel= zap.InfoLevel
	case "warn":
		log.LogLevel = zap.WarnLevel
	case "error":
		log.LogLevel = zap.ErrorLevel
	default:
		log.LogLevel = zap.DebugLevel
	}
}