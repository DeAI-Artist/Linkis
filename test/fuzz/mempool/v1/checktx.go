package v1

import (
	"github.com/DeAI-Artist/MintAI/abci/example/kvstore"
	"github.com/DeAI-Artist/MintAI/config"
	"github.com/DeAI-Artist/MintAI/libs/log"
	mempl "github.com/DeAI-Artist/MintAI/mempool"
	"github.com/DeAI-Artist/MintAI/proxy"

	mempoolv1 "github.com/DeAI-Artist/MintAI/mempool/v1"
)

var mempool mempl.Mempool

func init() {
	app := kvstore.NewApplication(config.GetDefaultDBDir())
	cc := proxy.NewLocalClientCreator(app)
	appConnMem, _ := cc.NewABCIClient()
	err := appConnMem.Start()
	if err != nil {
		panic(err)
	}
	cfg := config.DefaultMempoolConfig()
	cfg.Broadcast = false
	log := log.NewNopLogger()
	mempool = mempoolv1.NewTxMempool(log, cfg, appConnMem, 0)
}

func Fuzz(data []byte) int {

	err := mempool.CheckTx(data, nil, mempl.TxInfo{})
	if err != nil {
		return 0
	}

	return 1
}
