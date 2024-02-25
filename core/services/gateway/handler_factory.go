package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/DeAI-Artist/MintAI/core/chains/legacyevm"
	"github.com/DeAI-Artist/MintAI/core/logger"
	"github.com/DeAI-Artist/MintAI/core/services/gateway/config"
	"github.com/DeAI-Artist/MintAI/core/services/gateway/handlers"
	"github.com/DeAI-Artist/MintAI/core/services/gateway/handlers/functions"
	"github.com/DeAI-Artist/MintAI/core/services/pg"
)

const (
	FunctionsHandlerType HandlerType = "functions"
	DummyHandlerType     HandlerType = "dummy"
)

type handlerFactory struct {
	legacyChains legacyevm.LegacyChainContainer
	db           *sqlx.DB
	cfg          pg.QConfig
	lggr         logger.Logger
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, db *sqlx.DB, cfg pg.QConfig, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{legacyChains, db, cfg, lggr}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.db, hf.cfg, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
