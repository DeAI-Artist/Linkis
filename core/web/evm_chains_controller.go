package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

var ErrEVMNotEnabled = errChainDisabled{name: "EVM", tomlKey: "EVM.Enabled"}

func NewEVMChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.EVMChainResource](
		relay.EVM,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.EVM)),
		ErrEVMNotEnabled,
		presenters.NewEVMChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
