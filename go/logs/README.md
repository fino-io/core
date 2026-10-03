# Logs

`logs` 是当前的日志实现包，基于 `zap` 提供统一日志能力。

它提供两层 API：

- 包级默认 logger：适合大多数业务代码，直接调用 `logs.Info`、`logs.Infow`、`logs.Errorf`
- 独立 `Service`：适合需要独立 logger 实例的场景

支持能力：

- 控制台或文件输出
- 结构化字段日志
- 运行时动态调整日志级别
- 兼容 `server_logs` 作为旧包路径封装层

## 动态日志级别

当同时配置 `levelPattern` 和 `levelPort` 时，logger 会启动一个 HTTP 服务，暴露 `zap.AtomicLevel.ServeHTTP`：

- `GET <levelPattern>`：获取当前日志级别
- `PUT <levelPattern>`：动态修改日志级别

支持两种 `PUT` 请求体：

- `application/json`，例如 `{"level":"debug"}`
- `application/x-www-form-urlencoded`，例如 `level=debug`

## 主要接口

默认 logger：

```go
logger := logs.NewLoggerWith(&logs.Config{
  Level:  "info",
  Encode: "console",
  Output: "console",
})

logs.SetLogger(logger)
logs.SetLogLevel(logs.DebugLevel)
logs.Infow("service started", "name", "api")
```

独立 `Service`：

```go
svc := logs.NewService(logger)
svc.Infow("worker ready", "id", 7)
```

携带请求上下文：

```go
ctx = logs.WithFields(ctx,
  logs.Field{Key: "request_id", Value: requestID},
  logs.Field{Key: "tenant_id", Value: tenantID},
)
logs.Ctx(ctx).Infow("request handled", "status", "ok")
```

`WithFields` 只应附加允许进入日志的字段，不要放入 token、密码、请求体等敏感信息。
`Ctx` 会把 `ctx` 传到最底层 `Logger.Log`，便于上层适配 tracing 等能力；logger 不应保存
或异步复用传入的 context。

读取当前默认 logger：

```go
current := logs.DefaultLogger()
_ = current.GetLevel()
```

## 配置示例

```yaml
logs:
  level: info
  encode: console
  output: console
  levelPattern: /log/level
  levelPort: 22001
  file:
    path: ./log/app.log
    maxSize: 100
    maxBackups: 10
    maxAge: 30
    encode: json
```

`output: console` 输出到控制台；需要文件输出时设为 `file`。`file` 中的配置控制日志路径和轮转。

## 动态级别接口

如果配置了动态级别接口，可通过 HTTP 调整级别：

```bash
# 查看当前级别
curl http://127.0.0.1:22001/log/level

# 设置为 debug（JSON）
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"level":"debug"}' \
  http://127.0.0.1:22001/log/level

# 设置为 warn（表单）
curl -X PUT \
  -d "level=warn" \
  http://127.0.0.1:22001/log/level
```

说明：

- 创建带 `levelPattern` 和 `levelPort` 的 logger 时，会启动对应的 HTTP 服务
- 当前版本不会在 logger 替换时自动关闭旧的 level server
