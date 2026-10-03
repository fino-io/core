## `fino.core.File`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `name` | `string` |  | 否 |  |  |
| `isDir` | `boolean` |  | 否 |  |  |
| `mode` | [`fino.core.File.Mode`](file.md#fino-core-file-mode) |  | 否 |  |  |
| `info` | `fino.core.File.Info` |  | 否 |  |  |
| `files` | `Array<fino.core.File>` |  | 否 |  |  |


<a id="fino-core-file-mode"></a>

## `fino.core.File.Mode`

| 枚举值 | 数值 | 说明 |
|---|---:|---|
| `MODE_UNSPECIFIED` | 0 |  |
| `MODE_DIR` | 1 |  |
