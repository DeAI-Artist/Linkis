package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	abciclient "github.com/DeAI-Artist/Linkis/abci/client"
	"github.com/DeAI-Artist/Linkis/abci/example/kvstore"
	abciserver "github.com/DeAI-Artist/Linkis/abci/server"
	cfg "github.com/DeAI-Artist/Linkis/config"
)

func TestClientServerNoAddrPrefix(t *testing.T) {
	addr := "localhost:26658"
	transport := "socket"
	app := kvstore.NewApplication(cfg.GetDefaultDBDir())

	server, err := abciserver.NewServer(addr, transport, app)
	assert.NoError(t, err, "expected no error on NewServer")
	err = server.Start()
	assert.NoError(t, err, "expected no error on server.Start")

	client, err := abciclient.NewClient(addr, transport, true)
	assert.NoError(t, err, "expected no error on NewClient")
	err = client.Start()
	assert.NoError(t, err, "expected no error on client.Start")
}
