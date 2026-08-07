# fino Core Protobuf

Shared protobuf contracts for fino projects. The repository publishes one Buf
module and its Go runtime package.

| Module                  | Purpose | Go package                             |
|-------------------------| --- |----------------------------------------|
| `buf.build/fino/core`   | Common value types: time, duration, error, resource, file, URL, value, and version. | `github.com/fino-io/core/go/fino/core` |

## Use in a protobuf module

Add only the modules your schema imports:

```yaml
deps:
  - buf.build/fino/core
```

```proto
import "fino/core/time.proto";
import "fino2/options.proto";

message Book {
  option (fino2.model) = {generate: true services: "BookService"};

  fino.core.Timestamp create_time = 1;
}
```

Generated Go code imports the corresponding shared package. Do not vendor or
regenerate these proto files in individual services: a process must link each
protobuf descriptor only once.

## Maintain

Run Buf commands from the repository root:

```sh
buf lint
buf build
buf push
```

`buf push` publishes `fino/core`. Changing the protobuf package name, field
numbers, extension names, or extension field numbers is a public compatibility
change and must be treated as breaking.
