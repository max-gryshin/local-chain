package internal

import (
	"local-chain/internal/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMerkleTree_VerifyTransaction(t *testing.T) {
	tx1 := types.NewTransaction()
	tx1fake := types.NewTransaction()
	txs := []*types.Transaction{
		tx1,
		types.NewTransaction(),
		types.NewTransaction(),
		types.NewTransaction(),
		types.NewTransaction(),
	}
	tree, err := NewMerkleTree(txs)
	if err != nil {
		t.Error(err)
	}
	valid, err := tree.VerifyTransaction(tx1)
	if err != nil {
		t.Error(err)
	}
	require.True(t, valid)
	valid, err = tree.VerifyTransaction(tx1fake)
	if err != nil {
		require.Error(t, err)
	}
	require.False(t, valid, "Transaction should not be valid as it is a fake transaction")
}
