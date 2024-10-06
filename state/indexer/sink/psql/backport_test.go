package psql

import (
	"github.com/DeAI-Artist/Linkis/state/indexer"
	"github.com/DeAI-Artist/Linkis/state/txindex"
)

var (
	_ indexer.BlockIndexer = BackportBlockIndexer{}
	_ txindex.TxIndexer    = BackportTxIndexer{}
)
