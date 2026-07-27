# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [chaos/core/boxed.proto](#chaos_core_boxed-proto)
    - [BoolValue](#chaos-core-BoolValue)
    - [BoolValues](#chaos-core-BoolValues)
    - [BytesValue](#chaos-core-BytesValue)
    - [Float32Value](#chaos-core-Float32Value)
    - [Float32Values](#chaos-core-Float32Values)
    - [Float64Value](#chaos-core-Float64Value)
    - [Float64Values](#chaos-core-Float64Values)
    - [Int32Value](#chaos-core-Int32Value)
    - [Int32Values](#chaos-core-Int32Values)
    - [Int64Value](#chaos-core-Int64Value)
    - [Int64Values](#chaos-core-Int64Values)
    - [StringMap](#chaos-core-StringMap)
    - [StringMap.ValsEntry](#chaos-core-StringMap-ValsEntry)
    - [StringValue](#chaos-core-StringValue)
    - [StringValues](#chaos-core-StringValues)
    - [StringsMap](#chaos-core-StringsMap)
    - [StringsMap.ValsEntry](#chaos-core-StringsMap-ValsEntry)
    - [Uint32Value](#chaos-core-Uint32Value)
    - [Uint32Values](#chaos-core-Uint32Values)
    - [Uint64Value](#chaos-core-Uint64Value)
    - [Uint64Values](#chaos-core-Uint64Values)
  
- [chaos/core/duration.proto](#chaos_core_duration-proto)
    - [Duration](#chaos-core-Duration)
  
- [chaos/core/url.proto](#chaos_core_url-proto)
    - [Url](#chaos-core-Url)
    - [Url.Authority](#chaos-core-Url-Authority)
    - [Url.Query](#chaos-core-Url-Query)
    - [Url.Query.ValsEntry](#chaos-core-Url-Query-ValsEntry)
  
- [chaos/core/error_code.proto](#chaos_core_error_code-proto)
    - [ErrorCode](#chaos-core-ErrorCode)
  
- [chaos/core/null.proto](#chaos_core_null-proto)
    - [Null](#chaos-core-Null)
  
- [chaos/core/value.proto](#chaos_core_value-proto)
    - [Object](#chaos-core-Object)
    - [Object.ValsEntry](#chaos-core-Object-ValsEntry)
    - [Value](#chaos-core-Value)
    - [Values](#chaos-core-Values)
  
    - [ValueKind](#chaos-core-ValueKind)
  
- [chaos/core/error.proto](#chaos_core_error-proto)
    - [Error](#chaos-core-Error)
  
- [chaos/core/time.proto](#chaos_core_time-proto)
    - [Date](#chaos-core-Date)
    - [DateTime](#chaos-core-DateTime)
    - [TimeOfDay](#chaos-core-TimeOfDay)
    - [TimeZone](#chaos-core-TimeZone)
    - [Timestamp](#chaos-core-Timestamp)
  
    - [DayOfWeek](#chaos-core-DayOfWeek)
    - [Month](#chaos-core-Month)
  
- [chaos/core/file.proto](#chaos_core_file-proto)
    - [File](#chaos-core-File)
    - [File.Info](#chaos-core-File-Info)
  
    - [File.Mode](#chaos-core-File-Mode)
  
- [chaos/core/resource.proto](#chaos_core_resource-proto)
    - [Resource](#chaos-core-Resource)
  
- [chaos/core/version.proto](#chaos_core_version-proto)
    - [Version](#chaos-core-Version)
  
- [Scalar Value Types](#scalar-value-types)



<a name="chaos_core_boxed-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/boxed.proto



<a name="chaos-core-BoolValue"></a>

### BoolValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [bool](#bool) |  |  |






<a name="chaos-core-BoolValues"></a>

### BoolValues



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [bool](#bool) | repeated |  |






<a name="chaos-core-BytesValue"></a>

### BytesValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [bytes](#bytes) |  |  |






<a name="chaos-core-Float32Value"></a>

### Float32Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [float](#float) |  |  |






<a name="chaos-core-Float32Values"></a>

### Float32Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [float](#float) | repeated |  |






<a name="chaos-core-Float64Value"></a>

### Float64Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [double](#double) |  |  |






<a name="chaos-core-Float64Values"></a>

### Float64Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [double](#double) | repeated |  |






<a name="chaos-core-Int32Value"></a>

### Int32Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [int32](#int32) |  |  |






<a name="chaos-core-Int32Values"></a>

### Int32Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [int32](#int32) | repeated |  |






<a name="chaos-core-Int64Value"></a>

### Int64Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [int64](#int64) |  |  |






<a name="chaos-core-Int64Values"></a>

### Int64Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [int64](#int64) | repeated |  |






<a name="chaos-core-StringMap"></a>

### StringMap



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [StringMap.ValsEntry](#chaos-core-StringMap-ValsEntry) | repeated |  |






<a name="chaos-core-StringMap-ValsEntry"></a>

### StringMap.ValsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="chaos-core-StringValue"></a>

### StringValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [string](#string) |  |  |






<a name="chaos-core-StringValues"></a>

### StringValues



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [string](#string) | repeated |  |






<a name="chaos-core-StringsMap"></a>

### StringsMap



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [StringsMap.ValsEntry](#chaos-core-StringsMap-ValsEntry) | repeated |  |






<a name="chaos-core-StringsMap-ValsEntry"></a>

### StringsMap.ValsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [StringValues](#chaos-core-StringValues) |  |  |






<a name="chaos-core-Uint32Value"></a>

### Uint32Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [uint32](#uint32) |  |  |






<a name="chaos-core-Uint32Values"></a>

### Uint32Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [uint32](#uint32) | repeated |  |






<a name="chaos-core-Uint64Value"></a>

### Uint64Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| val | [uint64](#uint64) |  |  |






<a name="chaos-core-Uint64Values"></a>

### Uint64Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [uint64](#uint64) | repeated |  |





 

 

 

 



<a name="chaos_core_duration-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/duration.proto



<a name="chaos-core-Duration"></a>

### Duration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int64](#int64) |  |  |
| nanoseconds | [int32](#int32) |  |  |





 

 

 

 



<a name="chaos_core_url-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/url.proto



<a name="chaos-core-Url"></a>

### Url



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| scheme | [string](#string) |  |  |
| authority | [Url.Authority](#chaos-core-Url-Authority) |  |  |
| path | [string](#string) |  |  |
| query | [Url.Query](#chaos-core-Url-Query) |  |  |
| fragment | [string](#string) |  |  |






<a name="chaos-core-Url-Authority"></a>

### Url.Authority



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_info | [string](#string) |  |  |
| host | [string](#string) |  |  |
| port | [string](#string) |  |  |






<a name="chaos-core-Url-Query"></a>

### Url.Query



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [Url.Query.ValsEntry](#chaos-core-Url-Query-ValsEntry) | repeated |  |






<a name="chaos-core-Url-Query-ValsEntry"></a>

### Url.Query.ValsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [StringValues](#chaos-core-StringValues) |  |  |





 

 

 

 



<a name="chaos_core_error_code-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/error_code.proto



<a name="chaos-core-ErrorCode"></a>

### ErrorCode
ErrorCode
format: {domain}.{code}.{name}


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| name | [string](#string) |  | the name of error code |
| domain | [string](#string) |  | system, runtime, ... |
| description | [string](#string) |  | a detail description for the code |
| document | [Url](#chaos-core-Url) |  | the api document url for the error code |
| http_status_code | [int32](#int32) |  |  |





 

 

 

 



<a name="chaos_core_null-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/null.proto



<a name="chaos-core-Null"></a>

### Null






 

 

 

 



<a name="chaos_core_value-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/value.proto



<a name="chaos-core-Object"></a>

### Object



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [Object.ValsEntry](#chaos-core-Object-ValsEntry) | repeated |  |






<a name="chaos-core-Object-ValsEntry"></a>

### Object.ValsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Value](#chaos-core-Value) |  |  |






<a name="chaos-core-Value"></a>

### Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| null_value | [Null](#chaos-core-Null) |  |  |
| bool_value | [bool](#bool) |  |  |
| positive_value | [uint64](#uint64) |  |  |
| negative_value | [uint64](#uint64) |  |  |
| number_value | [double](#double) |  |  |
| string_value | [string](#string) |  |  |
| bytes_value | [bytes](#bytes) |  |  |
| object_value | [Object](#chaos-core-Object) |  |  |
| values_value | [Values](#chaos-core-Values) |  |  |






<a name="chaos-core-Values"></a>

### Values



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vals | [Value](#chaos-core-Value) | repeated |  |





 


<a name="chaos-core-ValueKind"></a>

### ValueKind


| Name | Number | Description |
| ---- | ------ | ----------- |
| VALUE_KIND_UNSPECIFIED | 0 |  |
| VALUE_KIND_NULL | 1 |  |
| VALUE_KIND_BOOLEAN | 2 |  |
| VALUE_KIND_INTEGER | 3 |  |
| VALUE_KIND_NUMBER | 4 |  |
| VALUE_KIND_STRING | 5 |  |
| VALUE_KIND_BYTES | 6 |  |
| VALUE_KIND_ARRAY | 7 |  |
| VALUE_KIND_OBJECT | 8 |  |


 

 

 



<a name="chaos_core_error-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/error.proto



<a name="chaos-core-Error"></a>

### Error



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [ErrorCode](#chaos-core-ErrorCode) |  |  |
| message | [string](#string) |  | a developer-facing error message |
| details | [Value](#chaos-core-Value) | repeated | a list if messages that carry the error details |





 

 

 

 



<a name="chaos_core_time-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/time.proto



<a name="chaos-core-Date"></a>

### Date



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| year | [int32](#int32) |  |  |
| month | [int32](#int32) |  |  |
| day | [int32](#int32) |  |  |






<a name="chaos-core-DateTime"></a>

### DateTime



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| year | [int32](#int32) |  |  |
| month | [int32](#int32) |  |  |
| day | [int32](#int32) |  |  |
| hour | [int32](#int32) |  |  |
| minute | [int32](#int32) |  |  |
| seconds | [int32](#int32) |  |  |
| nanoseconds | [int32](#int32) |  |  |
| time_zone | [TimeZone](#chaos-core-TimeZone) |  |  |






<a name="chaos-core-TimeOfDay"></a>

### TimeOfDay



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hours | [int32](#int32) |  |  |
| minutes | [int32](#int32) |  |  |
| seconds | [int32](#int32) |  |  |
| nanoseconds | [int32](#int32) |  |  |






<a name="chaos-core-TimeZone"></a>

### TimeZone



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| offset | [int32](#int32) |  |  |
| name | [string](#string) |  |  |






<a name="chaos-core-Timestamp"></a>

### Timestamp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int64](#int64) |  |  |
| nanoseconds | [int32](#int32) |  |  |





 


<a name="chaos-core-DayOfWeek"></a>

### DayOfWeek


| Name | Number | Description |
| ---- | ------ | ----------- |
| DAY_OF_WEEK_UNSPECIFIED | 0 |  |
| DAY_OF_WEEK_MONDAY | 1 |  |
| DAY_OF_WEEK_TUESDAY | 2 |  |
| DAY_OF_WEEK_WEDNESDAY | 3 |  |
| DAY_OF_WEEK_THURSDAY | 4 |  |
| DAY_OF_WEEK_FRIDAY | 5 |  |
| DAY_OF_WEEK_SATURDAY | 6 |  |
| DAY_OF_WEEK_SUNDAY | 7 |  |



<a name="chaos-core-Month"></a>

### Month


| Name | Number | Description |
| ---- | ------ | ----------- |
| MONTH_UNSPECIFIED | 0 |  |
| MONTH_JANUARY | 1 |  |
| MONTH_FEBRUARY | 2 |  |
| MONTH_MARCH | 3 |  |
| MONTH_APRIL | 4 |  |
| MONTH_MAY | 5 |  |
| MONTH_JUNE | 6 |  |
| MONTH_JULY | 7 |  |
| MONTH_AUGUST | 8 |  |
| MONTH_SEPTEMBER | 9 |  |
| MONTH_OCTOBER | 10 |  |
| MONTH_NOVEMBER | 11 |  |
| MONTH_DECEMBER | 12 |  |


 

 

 



<a name="chaos_core_file-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/file.proto



<a name="chaos-core-File"></a>

### File



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| is_dir | [bool](#bool) |  |  |
| mode | [File.Mode](#chaos-core-File-Mode) |  |  |
| info | [File.Info](#chaos-core-File-Info) |  |  |
| files | [File](#chaos-core-File) | repeated |  |






<a name="chaos-core-File-Info"></a>

### File.Info



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| suffix | [string](#string) |  |  |
| size | [int64](#int64) |  |  |
| change_time | [Timestamp](#chaos-core-Timestamp) |  |  |
| modify_time | [Timestamp](#chaos-core-Timestamp) |  |  |





 


<a name="chaos-core-File-Mode"></a>

### File.Mode


| Name | Number | Description |
| ---- | ------ | ----------- |
| MODE_UNSPECIFIED | 0 |  |
| MODE_DIR | 1 |  |


 

 

 



<a name="chaos_core_resource-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/resource.proto



<a name="chaos-core-Resource"></a>

### Resource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |





 

 

 

 



<a name="chaos_core_version-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaos/core/version.proto



<a name="chaos-core-Version"></a>

### Version



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| proto | [string](#string) |  |  |
| major | [int32](#int32) |  |  |
| minor | [int32](#int32) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

