## `fino.core.Timestamp`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `seconds` | `integer` | `Int64` | 否 |  |  |
| `nanoseconds` | `integer` | `Int32` | 否 |  |  |


## `fino.core.TimeZone`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `offset` | `integer` | `Int32` | 否 |  |  |
| `name` | `string` |  | 否 |  |  |


## `fino.core.TimeOfDay`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `hours` | `integer` | `Int32` | 否 |  |  |
| `minutes` | `integer` | `Int32` | 否 |  |  |
| `seconds` | `integer` | `Int32` | 否 |  |  |
| `nanoseconds` | `integer` | `Int32` | 否 |  |  |


## `fino.core.Date`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `year` | `integer` | `Int32` | 否 |  |  |
| `month` | `integer` | `Int32` | 否 |  |  |
| `day` | `integer` | `Int32` | 否 |  |  |


## `fino.core.DateTime`

| 字段 | 类型 | 格式 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|---|
| `year` | `integer` | `Int32` | 否 |  |  |
| `month` | `integer` | `Int32` | 否 |  |  |
| `day` | `integer` | `Int32` | 否 |  |  |
| `hour` | `integer` | `Int32` | 否 |  |  |
| `minute` | `integer` | `Int32` | 否 |  |  |
| `seconds` | `integer` | `Int32` | 否 |  |  |
| `nanoseconds` | `integer` | `Int32` | 否 |  |  |
| `timeZone` | `fino.core.TimeZone` |  | 否 |  |  |


<a id="fino-core-month"></a>

## `fino.core.Month`

| 枚举值 | 数值 | 说明 |
|---|---:|---|
| `MONTH_UNSPECIFIED` | 0 |  |
| `MONTH_JANUARY` | 1 |  |
| `MONTH_FEBRUARY` | 2 |  |
| `MONTH_MARCH` | 3 |  |
| `MONTH_APRIL` | 4 |  |
| `MONTH_MAY` | 5 |  |
| `MONTH_JUNE` | 6 |  |
| `MONTH_JULY` | 7 |  |
| `MONTH_AUGUST` | 8 |  |
| `MONTH_SEPTEMBER` | 9 |  |
| `MONTH_OCTOBER` | 10 |  |
| `MONTH_NOVEMBER` | 11 |  |
| `MONTH_DECEMBER` | 12 |  |


<a id="fino-core-dayofweek"></a>

## `fino.core.DayOfWeek`

| 枚举值 | 数值 | 说明 |
|---|---:|---|
| `DAY_OF_WEEK_UNSPECIFIED` | 0 |  |
| `DAY_OF_WEEK_MONDAY` | 1 |  |
| `DAY_OF_WEEK_TUESDAY` | 2 |  |
| `DAY_OF_WEEK_WEDNESDAY` | 3 |  |
| `DAY_OF_WEEK_THURSDAY` | 4 |  |
| `DAY_OF_WEEK_FRIDAY` | 5 |  |
| `DAY_OF_WEEK_SATURDAY` | 6 |  |
| `DAY_OF_WEEK_SUNDAY` | 7 |  |
