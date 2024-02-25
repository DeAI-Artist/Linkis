package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/keystore/keys/dkgsignkey"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewDKGSignKeysController(app chainlink.Application) KeysController {
	return NewKeysController[dkgsignkey.Key, presenters.DKGSignKeyResource](
		app.GetKeyStore().DKGSign(),
		app.GetLogger(),
		app.GetAuditLogger(),
		"dkgsignKey",
		presenters.NewDKGSignKeyResource,
		presenters.NewDKGSignKeyResources)
}
