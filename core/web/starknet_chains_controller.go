package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.StarkNet,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.StarkNet)),
		ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
