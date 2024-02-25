package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

// ErrCosmosNotEnabled is returned when COSMOS_ENABLED is not true.
var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", tomlKey: "Cosmos.Enabled"}

func NewCosmosNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.Cosmos)

	return newNodesController[presenters.CosmosNodeResource](
		scopedNodeStatuser, ErrCosmosNotEnabled, presenters.NewCosmosNodeResource, app.GetAuditLogger(),
	)
}
