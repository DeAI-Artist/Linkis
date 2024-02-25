package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/keystore/keys/solkey"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewSolanaKeysController(app chainlink.Application) KeysController {
	return NewKeysController[solkey.Key, presenters.SolanaKeyResource](app.GetKeyStore().Solana(), app.GetLogger(), app.GetAuditLogger(),
		"solanaKey", presenters.NewSolanaKeyResource, presenters.NewSolanaKeyResources)
}
