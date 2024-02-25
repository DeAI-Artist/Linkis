package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.CosmosChainResource](
		relay.Cosmos,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.Cosmos)),
		ErrCosmosNotEnabled,
		presenters.NewCosmosChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
