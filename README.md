# fino Core Protobuf

Shared protobuf contracts and Go runtime packages for fino projects.

| Go package | Purpose |
|------------|---------|
| `github.com/fino-io/core/go/fino/core` | Common value types and runtime helpers. |
| `github.com/fino-io/core/go/fino` | Code-generation options. |
| `github.com/fino-io/core/go/logs` | Shared logging runtime. |

See [ARCHITECTURE.md](ARCHITECTURE.md) for the function map, dependency
boundaries, runtime conventions, and known limitations.
See [CHANGELOG.md](CHANGELOG.md) for release changes and migration notes.

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
fino build --output=. --targets=api
```

API generation needs `package.yaml` and `proto/`. Service infrastructure
configuration in `fino.yaml` is only required when building the `service`
target. Use a Fino build containing the API-only configuration and
source-relative helper path fixes; the original v1.7.1 release requires
`fino.yaml` and writes this project's enum helpers to the wrong directory.

Generation updates:

- `go/fino/`: protobuf Go code and enum formatting/JSON helpers, kept in Git.
- `docs/`, `openapi/`, and `.fino/ir/`: local documentation, schemas, and
  descriptor artifacts, ignored by Git and regenerated on demand.

The hand-written value-type extensions, logging runtime, and their tests are
maintained alongside generated code in `go/`.

Validate the Go packages:

```sh
make -C go test
```

`go/Makefile` uses the Fino Go module template. Run its test and check targets
with `make -C go <target>`, including `test-fast`, `test-full`, `test-race`,
`test-coverage`, `vet`, `lint`, `sec`, `vuln`, and `verify`.

Lint uses the repository's `go/.golangci.yml`, including `staticcheck` and
`unused`. The GitHub workflow tests Go 1.22 and 1.25 and runs the repository's
lint configuration. Go tests verify that protobuf JSON tags match the message
descriptors. Generate protobuf files and enum helpers with Fino as described
above.

Changing the protobuf package name, field numbers, extension names, or extension
field numbers is a public compatibility change and must be treated as breaking.
