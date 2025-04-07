package aerospike

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aerospike/aerospike-client-go/v8"
	"github.com/stretchr/testify/assert"
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
	port, err := container.ServicePort(ctx)
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

func TestPutWithEnterprise(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainer(ctx, WithNamespace("namespace"), WithEnterpriseEdition())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoErrorf(t, err, "failed to terminate Aerospike container")
	})

	host, err := container.Host(ctx)
	require.NoErrorf(t, err, "failed to fetch Aerospike host")

	port, err := container.ServicePort(ctx)
	require.NoErrorf(t, err, "failed to fetch Aerospike port")

	client, err := aerospike.NewClient(host, port)
	require.NoErrorf(t, err, "failed to initialize Aerospike client")
	require.Truef(t, client.IsConnected(), "failed to connect to Aerospike")

	t.Run("Put", func(t *testing.T) {
		key, err := aerospike.NewKey("namespace", "set", "key1")
		require.NoErrorf(t, err, "failed to create Aerospike key")

		bin := aerospike.NewBin("bin", "value")

		err = client.PutBins(nil, key, bin)
		require.NoErrorf(t, err, "failed to create Aerospike record")
	})

	t.Run("Put with transaction and abort", func(t *testing.T) {
		// Begin a transaction
		policy := aerospike.NewWritePolicy(0, 0)

		txn := aerospike.NewTxn()
		policy.Txn = txn

		key, err := aerospike.NewKey("namespace", "set", "key2")
		require.NoErrorf(t, err, "failed to create Aerospike key")

		bin := aerospike.NewBin("bin", "value")

		_, err = client.Get(nil, key)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, aerospike.ErrKeyNotFound))

		// Perform put operations within the transaction
		err = client.PutBins(policy, key, bin)
		if err != nil && strings.Contains(err.Error(), "UNSUPPORTED_FEATURE") {
			t.Skip("Cluster is empty, skipping transaction test")
		}

		require.NoErrorf(t, err, "failed to create Aerospike record")

		// Verify that the record exists
		value, err := client.Get(nil, key)
		require.NoError(t, err)
		assert.Equal(t, "value", value.Bins["bin"])

		// Abort the transaction by not committing
		status, err := client.Abort(txn)
		require.NoError(t, err)
		assert.Equal(t, aerospike.AbortStatusOK, status)

		// Verify that the record does not exist
		_, err = client.Get(nil, key)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, aerospike.ErrKeyNotFound))
	})
}
