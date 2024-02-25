package web

import (
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
	"github.com/DeAI-Artist/MintAI/core/services/keystore/keys/cosmoskey"
	"github.com/DeAI-Artist/MintAI/core/web/presenters"
)

func NewCosmosKeysController(app chainlink.Application) KeysController {
	return NewKeysController[cosmoskey.Key, presenters.CosmosKeyResource](app.GetKeyStore().Cosmos(), app.GetLogger(), app.GetAuditLogger(),
		"cosmosKey", presenters.NewCosmosKeyResource, presenters.NewCosmosKeyResources)
}
