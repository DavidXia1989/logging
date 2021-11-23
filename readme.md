## Installation logging

zm日志 基于go.uber.org/zap/zapcore和github.com/lestrrat-go/file-rotatelogs的封装
Env：

golang >= 1.13.0

Install:

```
// 安装
go get -u github.com/DavidXia1989/logging
```

```
// 默认配置根目录下的log目录
// 默认日志文件名为app.log
// 默认每天切割日志，保存14天日志
```