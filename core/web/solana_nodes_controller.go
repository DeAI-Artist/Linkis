package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when Solana.Enabled is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", tomlKey: "Solana.Enabled"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.Solana)

	return newNodesController[presenters.SolanaNodeResource](
		scopedNodeStatuser, ErrSolanaNotEnabled, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}
