## `fino.ModelOptions`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `generate` | `boolean` |  | 否 |  |  |
| `services` | `Array<string>` |  | 否 |  |  |


## `fino.DBOptions`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `null` | `boolean` |  | 否 |  |  |
| `autoIncrement` | `boolean` |  | 否 |  |  |
| `len` | `integer` | `Int64` | 否 |  |  |
| `default` | `string` |  | 否 |  |  |
| `comment` | `string` |  | 否 |  |  |
| `index` | `string` |  | 否 |  |  |
| `uniqueIndex` | `string` |  | 否 |  |  |


## `fino.ValidateOptions`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `required` | `boolean` |  | 否 |  |  |
| `omitempty` | `boolean` |  | 否 |  |  |
| `dive` | `boolean` |  | 否 |  |  |
| `len` | `integer` | `Int64` | 否 |  |  |
| `min` | `integer` | `Int64` | 否 |  |  |
| `max` | `integer` | `Int64` | 否 |  |  |
| `oneof` | `string` |  | 否 |  |  |
| `email` | `boolean` |  | 否 |  |  |
| `url` | `boolean` |  | 否 |  |  |
| `uuid` | `boolean` |  | 否 |  |  |


## `fino.MessagingOptions`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `subscription` | `boolean` |  | 否 |  |  |
