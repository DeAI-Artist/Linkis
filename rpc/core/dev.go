package core

import (
	ctypes "github.com/DeAI-Artist/MintAI/rpc/core/types"
	rpctypes "github.com/DeAI-Artist/MintAI/rpc/jsonrpc/types"
)

// UnsafeFlushMempool removes all transactions from the mempool.
func UnsafeFlushMempool(ctx *rpctypes.Context) (*ctypes.ResultUnsafeFlushMempool, error) {
	env.Mempool.Flush()
	return &ctypes.ResultUnsafeFlushMempool{}, nil
}
