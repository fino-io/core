# 变更记录

## 未发布

### Protobuf 契约收敛

- `DBOptions` 的 `len/default/comment/index/unique_index`、`ValidateOptions` 的 `len/min/max/oneof` 改为 `optional`，保留显式零值和空字符串。生成的 Go 字段变为指针，调用方使用 `proto.Int64`、`proto.String` 赋值，nil 表示缺省。
- `File` 删除 `is_dir`，统一由 `mode` 表达类型；新增 `MODE_FILE`。`File.Info` 删除重复名称，统一使用 `File.name`。删除字段的名称和编号保留为 reserved，不提供旧字段适配。
- 补充时间单位、范围、时区优先级、Value 数值与紧凑 JSON 协议、URL 支持范围以及错误身份的契约注释。`change_time` 明确为元数据变更时间；修正 ErrorCode 格式与 Error.details 的旧注释。
- 使用包含共享枚举模板的 Fino 重建 API、文档、OpenAPI 和 IR，补充 options 存在性、字段选项扩展及文件树往返测试。

### 边界校验与基础能力补齐

- 查询 `Unmarshal` 将每个重复参数作为一个列表元素，保留逗号、空格、引号、空字符串和方括号文本。原来 `?tag=a,b` 或 `?tag=["a","b"]` 的隐式列表解析需改为重复参数，或显式调用 `UnmarshalParam`。固定数组要求元素数量一致。参数解码失败保留目标值。
- `Value`、`Object`、`Timestamp`、`Duration`、`Url` 支持标准 `encoding/json` 的紧凑协议。此前标准库输出内部字段结构的调用方需同步迁移；契约 JSON 继续使用 `protojson`。
- 新增 `CheckValid`：Timestamp 限定公元 1–9999 年及合法纳秒；Duration 保留 int64 秒范围但要求规范的纳秒；Value/Object/Values 检查负数幅值、UTF-8 和循环引用。越界 JSON 整数直接返回错误，浮点需明确包含小数点或指数。
- Date、TimeOfDay、DateTime 新增解析、格式化、校验及原生时间转换；TimeZone 支持 IANA 名称优先、固定 UTC 偏移和 nil 时区的明确规则。
- `logs.NewLoggerWith` 改为返回 `(Logger, error)`，调用方必须处理配置错误；新增 `NewLoggerWithWriter`，不接管 writer 的关闭。
- 补充核心回归/fuzz 用例，启用仓库内完整 lint、Go 版本矩阵和 protobuf JSON 标签契约测试。
- 复核补齐 jsoniter 循环值编码保护、查询指针列表的字面解析、被时区规则跳过的日期校验，以及自定义日志 writer 的空指针校验。

本轮收敛 Go 运行时的数据转换、错误分类、时间处理和日志生命周期。包含 Protobuf、Go API、紧凑 JSON 与 SQL 存储行为变化，调用方需要同步调整。

### API 调整

| 原用法 | 当前用法 |
| --- | --- |
| `Object.From2` | 使用 `Object.From`，成功替换全部字段并深复制，失败保留原值。 |
| 分类错误类型，例如 `*NotFoundError` | 构造器统一返回 `*Error`；通过 `IsNotFoundError` 等分类函数或 `errors.Is` 判断。 |
| `NewUnknownErrorError`、`IsUnknownErrorError`、`NewInternalErrorError` | 分别使用 `NewUnknownError`、`IsUnknownError`、`NewInternalError`。 |
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

- 普通 JSON 编解码直接使用 jsoniter 默认入口，移除额外全局配置；`Object.To` 保留完整浮点精度，动态数据转换中的数字通过 decoder 的 `UseNumber` 保留为 `json.Number`。
- 包装类型直接在外层 iterator/stream 上编解码，HTML 转义和映射键排序遵循调用方配置。集合注册继续使用 `RegisterJSONValuesCodec`，独立的 `ValsCodec` 已移除。
- Value 字符串与 Object 键在 JSON 编解码时拒绝非法 UTF-8，字段解码错误不替换原值。普通 JSON 输出遵循 jsoniter 默认 HTML 转义。
- 动态数据转换复用标准库完整 JSON 校验，非法 RawMessage 返回错误，不再静默替换为 null。包装类型的标准 JSON 方法直接处理字段，尾随内容导致的解码错误保留原值。
- `NewValue` 递归转换原生集合，保留字节、浮点和整数种类；映射键必须是字符串。结构体继续使用 JSON 标签和已注册的编码器。
- nil 数组输出 `null`，显式空切片输出 `[]`；无参数调用 `New*ArrayValue()` 等同于 nil 数组。数组读取保留这一形态。
- Value 浮点 JSON 保留小数点或指数，例如 `1.0`，避免解码成整数；普通保留字符串使用 `str.` 转义，二进制使用 `b64.`。
- 包装类型解码成功后替换字段；类型解码失败保留原值，null 清零或清空。`FromUrlValues` 替换全部键并复制源切片。
- Error JSON 的 `code` 从字符串改为完整 ErrorCode 对象，保留元数据。gokit 的扁平响应仍使用字符串 `code`，带 `error` 字段的响应需要按新对象格式读取 `error.code`。
- 枚举接受已知名称与严格整数，拒绝未知值和非法数字。
- Duration 的 SQL 输出改为精确十进制秒字符串，GORM 类型为 `text`。已有数值列需要由应用迁移；浮点列中已经丢失的精度无法恢复。
- Timestamp 输出 UTC RFC3339Nano，直接解析原始时区文本；URL 查询中的 `+` 由 `net/url` 编码为 `%2B`。
- 无 scheme 的 URL authority 保留 `//`，例如 `//example.com/path`；该结果可继续作为 URL 解析。

### 日志与生成

日志级别 handler 由应用挂载到已有 HTTP 服务；`LevelPort`、`LevelPattern` 已移除。派生 logger 共享级别和文件资源，在最后一个使用者停止后关闭；替换默认 logger 不自动关闭旧资源。

枚举 JSON helper 复用 `RegisterJSONEnum`，Fino 生成模板已同步调整。后续再生成时应使用包含该模板更新的 Fino 构建。

### 验证

在 `go/` 下执行 `go test ./...`、`go test -race ./...` 和 `go vet ./...`；在仓库根目录执行 `git diff --check`。发布版本号与 tag 由发布流程确定。
