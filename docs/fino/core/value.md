## `fino.core.Object`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `vals` | `Map<string, fino.core.Value>` |  | 否 |  |  |


## `fino.core.Values`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `vals` | `Array<fino.core.Value>` |  | 否 |  |  |


## `fino.core.Value`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `nullValue` | `fino.core.Null` |  | 否 |  |  |
| `boolValue` | `boolean` |  | 否 |  |  |
| `positiveValue` | `integer` | `Uint64` | 否 |  |  |
| `negativeValue` | `integer` | `Uint64` | 否 |  |  |
| `numberValue` | `number` | `Double` | 否 |  |  |
| `stringValue` | `string` |  | 否 |  |  |
| `bytesValue` | `string` | `Bytes` | 否 |  |  |
| `objectValue` | `fino.core.Object` |  | 否 |  |  |
| `valuesValue` | `fino.core.Values` |  | 否 |  |  |


<a id="fino-core-valuekind"></a>

## `fino.core.ValueKind`

| 枚举值 | 数值 | 说明 |
|---|---:|---|
| `VALUE_KIND_UNSPECIFIED` | 0 |  |
| `VALUE_KIND_NULL` | 1 |  |
| `VALUE_KIND_BOOLEAN` | 2 |  |
| `VALUE_KIND_INTEGER` | 3 |  |
| `VALUE_KIND_NUMBER` | 4 |  |
| `VALUE_KIND_STRING` | 5 |  |
| `VALUE_KIND_BYTES` | 6 |  |
| `VALUE_KIND_ARRAY` | 7 |  |
| `VALUE_KIND_OBJECT` | 8 |  |
