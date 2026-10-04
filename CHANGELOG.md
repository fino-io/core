# 变更记录

## 未发布

本轮收敛 Go 运行时的数据转换、错误分类、时间处理和日志生命周期。包含 Go API、紧凑 JSON 与 SQL 存储行为变化，调用方需要同步调整。Protobuf 契约及 `.pb.go` 保持不变。

### API 调整

| 原用法 | 当前用法 |
| --- | --- |
| `Object.From2` | 使用 `Object.From`，成功替换全部字段并深复制，失败保留原值。 |
| 分类错误类型，例如 `*NotFoundError` | 构造器统一返回 `*Error`；通过 `IsNotFoundError` 等分类函数或 `errors.Is` 判断。 |
| `NewDuration`、`Duration.FromSeconds` | 处理新增的 error 返回值；拒绝非有限或越界输入。 |
| `Duration.ToDuration`、`ToNanoSeconds`、`Timestamp.Add` | 处理新增的 error 返回值，转换溢出返回错误。 |
| `NewUrlQuery`、查询 `Add/Set` | 处理 error 返回值；输入只接受标量或一层标量集合。 |
| `StringValues.Contains` | 直接使用 bool 结果。 |
| `Value.GetValueArray` | 使用 `Value.GetValues`；`Object.GetValueArray(key)` 仍可用于读取指定字段。 |
| `Url.FormatWithoutSchema` | 改为 `FormatWithoutScheme`，作为省略 scheme 的展示文本使用。 |
| `StringLike`、`NewValueCodec`、`ValueCodec.DecodeAny`、`ValuesCodec` | 字符串转换使用 `ToStringConverter`；JSON 入口使用 jsoniter，数组 codec 复用 `RegisterJSONValuesCodec`。 |
| `Service.SetLogger` | 使用 `WithLogger` 返回派生服务；全局替换仍使用 `logs.SetLogger`。 |

`Error.AddDetail` 和对象键名转换也需要处理转换错误。`NewObjectFrom` 失败返回 `nil, error`。自定义 `logs.Logger` 实现需要提供 `LevelHandler`、`Sync` 和 `Close`。

### 数据行为

- `NewValue` 递归转换原生集合，保留字节、浮点和整数种类；映射键必须是字符串。结构体继续使用 JSON 标签和已注册的编码器。
- nil 数组输出 `null`，显式空切片输出 `[]`；无参数调用 `New*ArrayValue()` 等同于 nil 数组。数组读取保留这一形态。
- Value 浮点 JSON 保留小数点或指数，例如 `1.0`，避免解码成整数；普通保留字符串使用 `str.` 转义，二进制使用 `b64.`。
- 包装类型解码成功后替换字段；类型解码失败保留原值，null 清零或清空。`FromUrlValues` 替换全部键并复制源切片。
- Error JSON 保留完整 ErrorCode 元数据；枚举接受已知名称与严格整数，拒绝未知值和非法数字。
- Duration 的 SQL 输出改为精确十进制秒字符串，GORM 类型为 `text`。已有数值列需要由应用迁移；浮点列中已经丢失的精度无法恢复。
- Timestamp 输出 UTC RFC3339Nano，直接解析原始时区文本；URL 查询中的 `+` 由 `net/url` 编码为 `%2B`。
- 无 scheme 的 URL authority 保留 `//`，例如 `//example.com/path`；该结果可继续作为 URL 解析。

### 日志与生成

日志级别 handler 由应用挂载到已有 HTTP 服务；`LevelPort`、`LevelPattern` 已移除。派生 logger 共享级别和文件资源，在最后一个使用者停止后关闭；替换默认 logger 不自动关闭旧资源。

枚举 JSON helper 复用 `RegisterJSONEnum`，Fino 生成模板已同步调整。后续再生成时应使用包含该模板更新的 Fino 构建。

### 验证

在 `go/` 下执行 `go test ./...`、`go test -race ./...` 和 `go vet ./...`；在仓库根目录执行 `git diff --check`。发布版本号与 tag 由发布流程确定。
