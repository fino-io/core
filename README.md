# fino Core Protobuf

Shared protobuf contracts for fino projects. The repository publishes one Buf
module and the corresponding Go runtime packages.

| Module                | Purpose                                      | Go packages |
|-----------------------|----------------------------------------------|-------------|
| `buf.build/fino/core` | Common value types and code-generation options. | `github.com/fino-io/core/go/fino/core` (value types), `github.com/fino-io/core/go/fino` (options) |

## Use in a protobuf module

Add only the modules your schema imports:

```yaml
deps:
  - buf.build/fino/core
```

```proto
import "fino/core/time.proto";
import "fino/options.proto";

message Book {
  option (fino.model) = {generate: true services: "BookService"};

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

`buf push` publishes the module configured as `buf.build/fino/core`. The
registry module name does not determine protobuf package names: consumers
import `fino/options.proto` and use `(fino.model)`. Changing the protobuf
package name, field numbers, extension names, or extension field numbers is a
public compatibility change and must be treated as breaking.
