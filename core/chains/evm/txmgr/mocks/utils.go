package mocks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	txmgrmocks "github.com/DeAI-Artist/MintAI/common/txmgr/mocks"
	"github.com/DeAI-Artist/MintAI/core/chains/evm/gas"
	evmtypes "github.com/DeAI-Artist/MintAI/core/chains/evm/types"
)

type MockEvmTxManager = txmgrmocks.TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return txmgrmocks.NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee](t)
}
