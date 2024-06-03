package core

import (
	abci "github.com/DeAI-Artist/MintAI/abci/types"
	"github.com/DeAI-Artist/MintAI/libs/bytes"
	"github.com/DeAI-Artist/MintAI/proxy"
	ctypes "github.com/DeAI-Artist/MintAI/rpc/core/types"
	rpctypes "github.com/DeAI-Artist/MintAI/rpc/jsonrpc/types"
)

// ABCIQuery queries the application for some information.
// More: https://docs.tendermint.com/v0.34/rpc/#/ABCI/abci_query
func ABCIQuery(
	ctx *rpctypes.Context,
	path string,
	data bytes.HexBytes,
	height int64,
	prove bool,
) (*ctypes.ResultABCIQuery, error) {
	resQuery, err := env.ProxyAppQuery.QuerySync(abci.RequestQuery{
		Path:   path,
		Data:   data,
		Height: height,
		Prove:  prove,
	})
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultABCIQuery{Response: *resQuery}, nil
}

// ABCIInfo gets some info about the application.
// More: https://docs.tendermint.com/v0.34/rpc/#/ABCI/abci_info
func ABCIInfo(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {
	resInfo, err := env.ProxyAppQuery.InfoSync(proxy.RequestInfo)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultABCIInfo{Response: *resInfo}, nil
}
