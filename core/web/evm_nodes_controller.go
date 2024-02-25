package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.EVM)

	return newNodesController[presenters.EVMNodeResource](
		scopedNodeStatuser, ErrEVMNotEnabled, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
