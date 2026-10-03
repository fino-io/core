ErrorCode
 format: {domain}.{code}.{name}

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `code` | `integer` | `Int32` | 否 |  |  |
| `name` | `string` |  | 否 |  | the name of error code |
| `domain` | `string` |  | 否 |  | system, runtime, ... |
| `description` | `string` |  | 否 |  | a detail description for the code |
| `document` | `string` | `Url` | 否 |  | the api document url for the error code |
| `httpStatusCode` | `integer` | `Int32` | 否 |  |  |
