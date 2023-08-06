# testcontainers-aerospike-go

Go library for [Aerospike](https://aerospike.com/) integration testing via
[Testcontainers](https://testcontainers.com/).

## Example

```go
import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    aero "github.com/aerospike/aerospike-client-go/v6"
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
    port, err := container.Port(ctx)
    require.NoError(t, err)

    client, err := aero.NewClient(host, port)
    require.NoError(t, err)

    return client
}
```
