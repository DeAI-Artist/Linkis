package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DeAI-Artist/MintAI/core/cmd"
	"github.com/DeAI-Artist/MintAI/core/internal/cltest"
	"github.com/DeAI-Artist/MintAI/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
)

func TestShell_IndexSolanaChains(t *testing.T) {
	t.Parallel()

	id := solanatest.RandomChainID()
	cfg := solana.TOMLConfig{
		ChainID: &id,
		Enabled: ptr(true),
	}
	app := solanaStartNewApplication(t, &cfg)
	client, r := app.NewShellAndRenderer()

	require.Nil(t, cmd.SolanaChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.SolanaChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, id, c.ID)
	assertTableRenders(t, r)
}
