package txmgr

import (
	"context"

	"github.com/DeAI-Artist/MintAI/common/types"
)

type SequenceSyncer[ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable, SEQ types.Sequence] interface {
	Sync(ctx context.Context, addr ADDR, localSequence SEQ) (SEQ, error)
}
