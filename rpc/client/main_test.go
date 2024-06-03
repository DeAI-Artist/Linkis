package client_test

import (
	"os"
	"testing"

	"github.com/DeAI-Artist/MintAI/abci/example/kvstore"
	nm "github.com/DeAI-Artist/MintAI/node"
	rpctest "github.com/DeAI-Artist/MintAI/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a tendermint node (and kvstore) in the background to test against
	dir, err := os.MkdirTemp("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}

	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartTendermint(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
