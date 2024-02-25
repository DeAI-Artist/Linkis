package headtracker

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/DeAI-Artist/MintAI/common/headtracker"
	commontypes "github.com/DeAI-Artist/MintAI/common/types"
	evmtypes "github.com/DeAI-Artist/MintAI/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type headBroadcaster = headtracker.HeadBroadcaster[*evmtypes.Head, common.Hash]

var _ commontypes.HeadBroadcaster[*evmtypes.Head, common.Hash] = &headBroadcaster{}

func NewHeadBroadcaster(
	lggr logger.Logger,
) *headBroadcaster {
	return headtracker.NewHeadBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
