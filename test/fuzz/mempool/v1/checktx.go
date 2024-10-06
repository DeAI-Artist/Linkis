package v1

import (
	"github.com/DeAI-Artist/Linkis/abci/example/kvstore"
	"github.com/DeAI-Artist/Linkis/config"
	"github.com/DeAI-Artist/Linkis/libs/log"
	mempl "github.com/DeAI-Artist/Linkis/mempool"
	"github.com/DeAI-Artist/Linkis/proxy"

	mempoolv1 "github.com/DeAI-Artist/Linkis/mempool/v1"
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
