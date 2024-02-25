package legacyevm

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	evmclient "github.com/DeAI-Artist/MintAI/core/chains/evm/client"
	evmconfig "github.com/DeAI-Artist/MintAI/core/chains/evm/config"
	"github.com/DeAI-Artist/MintAI/core/chains/evm/gas"
	"github.com/DeAI-Artist/MintAI/core/chains/evm/logpoller"
	"github.com/DeAI-Artist/MintAI/core/chains/evm/txmgr"
	"github.com/DeAI-Artist/MintAI/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.EVM,
	evmRPCEnabled bool,
	databaseConfig txmgr.DatabaseConfig,
	listenerConfig txmgr.ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainRelayExtenderConfig,
) (txm txmgr.TxManager,
	estimator gas.EvmFeeEstimator,
	err error,
) {
	chainID := cfg.ChainID()
	if !evmRPCEnabled {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
		return txm, nil, nil
	}

	lggr = lggr.Named("Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", cfg.GasEstimator().BumpTxDepth(),
		"maxInFlightTransactions", cfg.Transactions().MaxInFlight(),
		"maxQueuedTransactions", cfg.Transactions().MaxQueued(),
		"nonceAutoSync", cfg.NonceAutoSync(),
		"limitDefault", cfg.GasEstimator().LimitDefault(),
	)

	// build estimator from factory
	if opts.GenGasEstimator == nil {
		estimator = gas.NewEstimator(lggr, client, cfg, cfg.GasEstimator())
	} else {
		estimator = opts.GenGasEstimator(chainID)
	}

	if opts.GenTxManager == nil {
		txm, err = txmgr.NewTxm(
			db,
			cfg,
			txmgr.NewEvmTxmFeeConfig(cfg.GasEstimator()),
			cfg.Transactions(),
			databaseConfig,
			listenerConfig,
			client,
			lggr,
			logPoller,
			opts.KeyStore,
			estimator)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return
}
