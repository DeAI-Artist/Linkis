package v0

import (
	"github.com/DeAI-Artist/MintAI/abci/example/kvstore"
	"github.com/DeAI-Artist/MintAI/config"
	mempl "github.com/DeAI-Artist/MintAI/mempool"
	mempoolv0 "github.com/DeAI-Artist/MintAI/mempool/v0"
	"github.com/DeAI-Artist/MintAI/proxy"
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
	mempool = mempoolv0.NewCListMempool(cfg, appConnMem, 0)
}

func Fuzz(data []byte) int {
	err := mempool.CheckTx(data, nil, mempl.TxInfo{})
	if err != nil {
		return 0
	}

	return 1
}
