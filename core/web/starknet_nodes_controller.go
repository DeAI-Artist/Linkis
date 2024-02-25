package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

// ErrStarkNetNotEnabled is returned when Starknet.Enabled is not true.
var ErrStarkNetNotEnabled = errChainDisabled{name: "StarkNet", tomlKey: "Starknet.Enabled"}

func NewStarkNetNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.StarkNet)

	return newNodesController[presenters.StarkNetNodeResource](
		scopedNodeStatuser, ErrStarkNetNotEnabled, presenters.NewStarkNetNodeResource, app.GetAuditLogger())
}
