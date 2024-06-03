package null

import (
	"context"
	"errors"

	abci "github.com/DeAI-Artist/MintAI/abci/types"
	"github.com/DeAI-Artist/MintAI/libs/pubsub/query"
	"github.com/DeAI-Artist/MintAI/state/txindex"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

// TxIndex acts as a /dev/null.
type TxIndex struct{}

// Get on a TxIndex is disabled and panics when invoked.
func (txi *TxIndex) Get(hash []byte) (*abci.TxResult, error) {
	return nil, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)
}

// AddBatch is a noop and always returns nil.
func (txi *TxIndex) AddBatch(batch *txindex.Batch) error {
	return nil
}

// Index is a noop and always returns nil.
func (txi *TxIndex) Index(result *abci.TxResult) error {
	return nil
}

func (txi *TxIndex) Search(ctx context.Context, q *query.Query) ([]*abci.TxResult, error) {
	return []*abci.TxResult{}, nil
}
