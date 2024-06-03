package psql

import (
	"github.com/DeAI-Artist/MintAI/state/indexer"
	"github.com/DeAI-Artist/MintAI/state/txindex"
)

var (
	_ indexer.BlockIndexer = BackportBlockIndexer{}
	_ txindex.TxIndexer    = BackportTxIndexer{}
)
