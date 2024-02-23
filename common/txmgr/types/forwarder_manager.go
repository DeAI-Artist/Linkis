package types

import (
	"github.com/DeAI-Artist/MintAI/common/types"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

//go:generate mockery --quiet --name ForwarderManager --output ./mocks/ --case=underscore
type ForwarderManager[ADDR types.Hashable] interface {
	services.Service
	ForwarderFor(addr ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
