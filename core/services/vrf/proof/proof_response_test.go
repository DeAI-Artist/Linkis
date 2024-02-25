package proof_test

import (
	"math/big"
	"testing"

	"github.com/DeAI-Artist/MintAI/core/gethwrappers/generated/solidity_vrf_verifier_wrapper"
	"github.com/DeAI-Artist/MintAI/core/internal/testutils/configtest"
	proof2 "github.com/DeAI-Artist/MintAI/core/services/vrf/proof"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/DeAI-Artist/MintAI/core/chains/evm/assets"
	"github.com/DeAI-Artist/MintAI/core/internal/cltest"
	"github.com/DeAI-Artist/MintAI/core/internal/testutils/pgtest"
)

func TestMarshaledProof(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	key := cltest.DefaultVRFKey
	require.NoError(t, keyStore.VRF().Add(key))
	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(1)
	s := proof2.TestXXXSeedData(t, preSeed, blockHash, blockNum)
	proofResponse, err := proof2.GenerateProofResponse(keyStore.VRF(), key.ID(), s)
	require.NoError(t, err)
	goProof, err := proof2.UnmarshalProofResponse(proofResponse)
	require.NoError(t, err)
	actualProof, err := goProof.CryptoProof(s)
	require.NoError(t, err)
	proof, err := proof2.MarshalForSolidityVerifier(&actualProof)
	require.NoError(t, err)
	// NB: For changes to the VRF solidity code to be reflected here, "go generate"
	// must be run in core/services/vrf.
	ethereumKey, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(ethereumKey, big.NewInt(1337))
	require.NoError(t, err)
	genesisData := core.GenesisAlloc{auth.From: {Balance: assets.Ether(100).ToInt()}}
	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	_, _, verifier, err := solidity_vrf_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
	_, err = verifier.RandomValueFromVRFProof(nil, proof[:])
	require.NoError(t, err)
}
