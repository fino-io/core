# fino Core Protobuf

Shared protobuf contracts and Go runtime packages for fino projects.

| Go package | Purpose |
|------------|---------|
| `github.com/fino-io/core/go/fino/core` | Common value types and runtime helpers. |
| `github.com/fino-io/core/go/fino` | Code-generation options. |

## Use in a Fino project

```proto
import "fino/core/time.proto";
import "fino/options.proto";

message Book {
  option (fino.model) = {generate: true services: "BookService"};

  fino.core.Timestamp create_time = 1;
}
```

Fino resolves these shared proto imports from its bundled runtime descriptors.
Generated Go code imports the corresponding shared package. Do not vendor or
regenerate these proto files in individual services: a process must link each
protobuf descriptor only once.

## Maintain

Generate from the protobuf contracts with Fino from the repository root:

```sh
fino doctor --output=. --targets=api
make build
```

`make build` runs `fino build --output=. --targets=api` and places the enum
helpers alongside their protobuf types in `go/fino/core`. This directory step
is needed with Fino v1.7.1, which writes those helpers to `go/core`.
The API-only configuration is recorded in `fino.yaml` with an empty backend.

Generation updates:

- `go/fino/`: protobuf Go code and enum formatting/JSON helpers.
- `docs/`: Markdown documentation for the protobuf contracts.
- `openapi/`: JSON schemas for the shared messages.
- `.fino/ir/`: local descriptor artifacts, ignored by Git.

Validate the Go packages:

```sh
make test
```

`go/Makefile` uses the Fino Go module template. The root Makefile forwards its
test and check targets to that module: `test-fast`, `test-full`, `test-race`,
`test-coverage`, `vet`, `lint`, `sec`, `vuln`, and `verify`.

Changing the protobuf package name, field numbers, extension names, or extension
field numbers is a public compatibility change and must be treated as breaking.
