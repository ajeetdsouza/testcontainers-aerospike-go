package aerospike

import (
	"context"
	"testing"

	"github.com/aerospike/aerospike-client-go/v6"
	"github.com/stretchr/testify/require"
)

func TestPut(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainer(ctx, WithNamespace("namespace"))
	require.NoError(t, err)
	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoErrorf(t, err, "failed to terminate Aerospike container")
	})

	host, err := container.Host(ctx)
	require.NoErrorf(t, err, "failed to fetch Aerospike host")
	port, err := container.Port(ctx)
	require.NoErrorf(t, err, "failed to fetch Aerospike port")

	client, err := aerospike.NewClient(host, port)
	require.NoErrorf(t, err, "failed to initialize Aerospike client")
	require.Truef(t, client.IsConnected(), "failed to connect to Aerospike")

	key, err := aerospike.NewKey("namespace", "set", "key")
	require.NoErrorf(t, err, "failed to create Aerospike key")
	bin := aerospike.NewBin("bin", "value")

	err = client.PutBins(nil, key, bin)
	require.NoErrorf(t, err, "failed to create Aerospike record")
}
