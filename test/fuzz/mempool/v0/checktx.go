package v0

import (
	"github.com/DeAI-Artist/Linkis/abci/example/kvstore"
	"github.com/DeAI-Artist/Linkis/config"
	mempl "github.com/DeAI-Artist/Linkis/mempool"
	mempoolv0 "github.com/DeAI-Artist/Linkis/mempool/v0"
	"github.com/DeAI-Artist/Linkis/proxy"
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
