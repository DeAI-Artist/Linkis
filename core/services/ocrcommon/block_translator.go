package ocrcommon

import (
	"context"
	"math/big"

	"github.com/DeAI-Artist/MintAI/common/config"
	evmclient "github.com/DeAI-Artist/MintAI/core/chains/evm/client"
	evmtypes "github.com/DeAI-Artist/MintAI/core/chains/evm/types"
	"github.com/DeAI-Artist/MintAI/core/logger"
)

// BlockTranslator converts emitted block numbers (from block.number) into a
// block number range suitable for query in FilterLogs
type BlockTranslator interface {
	NumberToQueryRange(ctx context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int)
}

// NewBlockTranslator returns the block translator for the given chain
func NewBlockTranslator(cfg Config, client evmclient.Client, lggr logger.Logger) BlockTranslator {
	switch cfg.ChainType() {
	case config.ChainArbitrum:
		return NewArbitrumBlockTranslator(client, lggr)
	case config.ChainXDai, config.ChainMetis, config.ChainOptimismBedrock:
		fallthrough
	default:
		return &l1BlockTranslator{}
	}
}

type l1BlockTranslator struct{}

func (*l1BlockTranslator) NumberToQueryRange(_ context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	return big.NewInt(int64(changedInL1Block)), big.NewInt(int64(changedInL1Block))
}

func (*l1BlockTranslator) OnNewLongestChain(context.Context, *evmtypes.Head) {}
