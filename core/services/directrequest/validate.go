package directrequest

import (
	"github.com/pkg/errors"

	"github.com/DeAI-Artist/MintAI/core/chains/evm/utils/big"
	"github.com/DeAI-Artist/MintAI/core/null"
	"github.com/DeAI-Artist/MintAI/core/services/job"
	"github.com/DeAI-Artist/MintAI/core/services/keystore/keys/ethkey"
	"github.com/DeAI-Artist/MintAI/core/store/models"
	"github.com/smartcontractkit/chainlink-common/pkg/assets"
)

type DirectRequestToml struct {
	ContractAddress          ethkey.EIP55Address      `toml:"contractAddress"`
	Requesters               models.AddressCollection `toml:"requesters"`
	MinContractPayment       *assets.Link             `toml:"minContractPaymentLinkJuels"`
	EVMChainID               *big.Big                 `toml:"evmChainID"`
	MinIncomingConfirmations null.Uint32              `toml:"minIncomingConfirmations"`
}

func ValidatedDirectRequestSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{}
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}
	var spec DirectRequestToml
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}
	jb.DirectRequestSpec = &job.DirectRequestSpec{
		ContractAddress:          spec.ContractAddress,
		Requesters:               spec.Requesters,
		MinContractPayment:       spec.MinContractPayment,
		EVMChainID:               spec.EVMChainID,
		MinIncomingConfirmations: spec.MinIncomingConfirmations,
	}

	if jb.Type != job.DirectRequest {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	return jb, nil
}
