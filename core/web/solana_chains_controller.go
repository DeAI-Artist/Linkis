package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/relay"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.Solana,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.Solana)),
		ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
