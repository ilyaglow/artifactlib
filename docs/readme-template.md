<h1 align="center">{{.Name}}</h1>

<p  align="center">
 <a href="https://{{.ModulePath}}/actions"><img src="https://{{.ModulePath}}/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/{{.RelModulePath}}"><img src="https://codecov.io/gh/{{.RelModulePath}}/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/{{.ModulePath}}"><img src="https://goreportcard.com/badge/{{.ModulePath}}" alt="report" /></a>
 <a href="https://pkg.go.dev/{{.ModulePath}}"><img src="https://godoc.org/{{.ModulePath}}?status.svg" alt="doc" /></a>
</p>

{{.Doc}}


## Python library

### Installation

Python installation can be easily done via pip:

```bash
pip install pyartifacts
```

### Usage

```python
import forensicstore

if __name__ == '__main__':
    registry = Registry()
    registry.read_folder("test/artifacts/valid")
    print(registry)
```

The full documentation can be found in [the documentation](TODO).

## Go package

### Installation


```bash
go get -u {{.ModulePath}}
```

{{if .Examples}}
### Usage
{{ range $key, $value := .Examples }}

{{if $key}}### {{ $key }}{{end}}
```go
{{ $value }}
```
{{end}}{{end}}

The full documentation can be found in the [godocs](https://pkg.go.dev/{{.ModulePath}}).