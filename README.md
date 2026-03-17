# structpolicy

Go static analyzers for enforcing struct usage policy in Go APIs. Built on [`golang.org/x/tools/go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis).

## Analyzers

| Analyzer | Policy | Details |
|----------|--------|---------|
| [ptrstruct](cmd/ptrstruct/) | Struct types must be used by **pointer** | Reports value usage of struct types |
| [valuestruct](cmd/valuestruct/) | Struct types must be used by **value** | Reports pointer usage of struct types |

The analyzers ship with mode-specific performance-tuning defaults. `ptrstruct` defaults focus on surfacing `copy hotspot` candidates, while `valuestruct` defaults focus on surfacing `allocation / indirection hotspot` candidates. See each analyzer's README for the exact default profile and full flag reference.

## Local Development

Requires [Task](https://taskfile.dev/) as the build interface.

```bash
task          # list available tasks
task build    # build the binary
task test     # run tests with race detection, shuffle, count=10
task lint     # run golangci-lint
task check    # full verification (lint + build + test + tidy)
```

## License

[MIT](LICENSE)
