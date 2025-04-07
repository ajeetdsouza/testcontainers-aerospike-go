<!-- markdownlint-configure-file {
  "MD033": false,
  "MD041": false
} -->

<div align="center">

# testcontainers-aerospike-go

[![Go Reference](https://pkg.go.dev/badge/github.com/ajeetdsouza/testcontainers-aerospike-go.svg)](https://pkg.go.dev/github.com/ajeetdsouza/testcontainers-aerospike-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajeetdsouza/testcontainers-aerospike-go)](https://goreportcard.com/report/github.com/ajeetdsouza/testcontainers-aerospike-go)

Go library for **[Aerospike](https://aerospike.com/) integration testing via
[Testcontainers](https://testcontainers.com/)**.

</div>

## Install

Use `go get` to install the latest version of the library.

```bash
go get -u github.com/ajeetdsouza/testcontainers-aerospike-go@latest
```

## Usage

```go
import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    aero "github.com/aerospike/aerospike-client-go/v8"
    aeroTest "github.com/ajeetdsouza/testcontainers-aerospike-go"
)

func TestAerospike(t *testing.T) {
    aeroClient := setupAerospike(t)
    // your code here
}

func setupAerospike(t *testing.T) *aero.Client {
    ctx := context.Background()

    container, err := aeroTest.RunContainer(ctx)
    require.NoError(t, err)
    t.Cleanup(func() {
        err := container.Terminate(ctx)
        require.NoError(t, err)
    })

    host, err := container.Host(ctx)
    require.NoError(t, err)
    port, err := container.ServicePort(ctx)
    require.NoError(t, err)

    client, err := aero.NewClient(host, port)
    require.NoError(t, err)

    return client
}
```
