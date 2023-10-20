package benchmarkor

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcClient struct {
	client      *ethclient.Client
	blockNumber *big.Int
}

type RpcCall func(client *ethclient.Client, callOpts *bind.CallOpts) (status int, err error)

func NewRpcClient(rpcUrl string) (rpcClient *RpcClient, err error) {
	timeout := time.Second * 10
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	rpcClient = &RpcClient{
		client: client,
	}
	return rpcClient, nil
}

func (c *RpcClient) rpcCall(rpcCallFn RpcCall, callOpts *bind.CallOpts) (status int, timeTaken uint64, err error) {
	start := time.Now()
	statusReturn, err := rpcCallFn(c.client, callOpts)
	diff := uint64(time.Since(start).Nanoseconds() / 1000)
	if err != nil {
		return statusReturn, diff, err
	}
	return statusReturn, diff, nil

}
