package apis

import (
	"benchmarkor-go/abigen/UniswapV3Pool"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/exp/slog"
)

func GetBlockByNumber(client *ethclient.Client, callOpts *bind.CallOpts) (status int, err error) {
	res, err := client.BlockByNumber(callOpts.Context, callOpts.BlockNumber)
	if err != nil {
		slog.Error("GetBlockByNumber", "err", err)
		return http.StatusBadRequest, err
	}
	slog.Info("Result", "res", res.Number())

	return http.StatusAccepted, nil
}

func GetBalanceAt(client *ethclient.Client, callOpts *bind.CallOpts) (status int, err error) {
	const VITALIK_ADDRESS_STRING = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	vitalik := common.HexToAddress(VITALIK_ADDRESS_STRING)

	res, err := client.BalanceAt(callOpts.Context, vitalik, callOpts.BlockNumber)
	if err != nil {
		slog.Error("GetBalanceAt", "err", err)
		return http.StatusBadRequest, err
	}
	slog.Info("Result", "res", res.String())

	return http.StatusAccepted, nil
}

func GetContractSlot0(client *ethclient.Client, callOpts *bind.CallOpts) (status int, err error) {
	const UNISWAP_ETH_USDC_POOL = "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"
	pool := common.HexToAddress(UNISWAP_ETH_USDC_POOL)

	uniswapContract, err := UniswapV3Pool.NewUniswapV3Pool(pool, client)
	if err != nil {
		slog.Warn("Unable to instantiate contract")
		return http.StatusBadRequest, err
	}

	res, err := uniswapContract.Slot0(callOpts)
	if err != nil {
		slog.Error("GetContractSlot0", "err", err)
		return http.StatusBadRequest, err
	}
	slog.Info("Result", "Res", res.SqrtPriceX96.String())

	return http.StatusAccepted, nil
}
