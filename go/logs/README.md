# Logs

日志运行时复用 Zap 和 Lumberjack，提供包级默认 logger、独立 Service、上下文字段、文件轮转及动态级别 handler。

## 使用

```go
logger, err := logs.NewLoggerWith(&logs.Config{
    Level: "info",
    Encode: "json",
    Output: "console",
})
if err != nil {
    return err
}
defer logger.Close()

logs.SetLogger(logger)
logs.Infow("service started", "name", "api")
logs.SetLogLevel(logs.DebugLevel)

svc := logs.NewService(logger)
svc.Infow("worker ready", "id", 7)
```

`Logger` 提供 `SetLevel`、`GetLevel`、`With`、`Log`、`LevelHandler`、`Sync`、`Close`。

`NewLoggerWith` 返回 `(Logger, error)`，拒绝非法级别、编码、输出类型及负数文件限制。
未指定的级别、编码、输出和文件路径使用默认值；文件保留数量和天数为零表示不限制。
配置初始化不会修改传入的 `Config`。

自定义输出使用 `NewLoggerWithWriter(cfg, writer)`，同样返回 `(Logger, error)`。
writer 由调用方管理，`Sync` 会转发同步操作，`Close` 不关闭调用方的 writer。

`Service` 的配置不可变。`WithLogger` 和 `WithContext` 返回新实例，父实例保持原配置；全局 `SetLogger` 可以并发替换默认 logger。派生 logger 共享级别和输出资源。

```go
ctx := logs.WithFields(ctx,
    logs.Field{Key: "request_id", Value: requestID},
)
logs.Ctx(ctx).Infow("request handled", "status", "ok")

child := svc.WithLogger(otherLogger)
child.Info("using another logger")
```

字段按初始配置、`Logger.With`、上下文、单次调用的顺序覆盖，同名字段只输出一次。`Ctx` 将 context 传到 `Logger.Log`。

## 动态级别

将 handler 挂载到应用已有的 HTTP 服务：

```go
mux.Handle("/log/level", logger.LevelHandler())
```

handler 直接复用 `zap.AtomicLevel`：GET 读取级别，PUT 更新级别，支持 JSON `{"level":"debug"}` 和表单 `level=debug`。服务的绑定地址、中间件、超时和关闭过程由应用管理。日志配置不再包含 `levelPort`、`levelPattern`，创建 logger 不启动 HTTP 服务。

## 文件与关闭

```yaml
logs:
  level: info
  encode: console
  output: file
  file:
    path: ./log/app.log
    maxSize: 100
    maxBackups: 10
    maxAge: 30
    encode: json
```

文件未指定编码时使用 JSON。`Sync` 转发底层输出的同步操作，`Close` 关闭 logger 拥有的轮转文件，并且可以重复调用；它不关闭控制台或调用方提供的 writer。Lumberjack 直接写文件，不需要额外缓冲层。

派生 logger 共享文件资源，应在最后一个使用者停止记录后关闭。替换全局 logger 不自动关闭旧实例，关闭时机由应用明确管理。
