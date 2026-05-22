# Module path

This module is published as:

```text
github.com/TiwariAshuism/apigen
```

The `go.mod` `module` line **must** match the GitHub repository path. If they differ, `go install` / `go run ...@latest` fails with a version constraints conflict.

## Use in other projects

```bash
go install github.com/TiwariAshuism/apigen/cmd/apigen@latest
```

```go
//go:generate go run github.com/TiwariAshuism/apigen/cmd/apigen@latest -input routes.go -output ..
```

Requires a tagged release or the latest commit on the default branch visible to the Go module proxy (`GOPROXY`).

## Local development (before push)

```bash
go run ./cmd/apigen -input api/routes.go -output /path/to/project -module your/module/path
```

Or in the consumer `go.mod`:

```go
replace github.com/TiwariAshuism/apigen => ../apigen
```
