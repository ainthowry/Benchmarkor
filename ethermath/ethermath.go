package ethermath

import (
	"math"
	"math/big"
)

const X96_String = "79228162514264337593543950336" //2**96 -> new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)

func SqrtPriceX96ToPrice(sqrtPriceX96 *big.Int, token0Decimal int, token1Decimal int) *big.Float {
	//price = (sqrtPriceX96 / (2 ** 96)) ** 2
	X96, _ := new(big.Int).SetString(X96_String, 10)
	price := new(big.Int)
	price = price.Div(sqrtPriceX96, X96).Exp(price, big.NewInt(2), nil)
	priceInDecimals := new(big.Float).SetInt(price)
	ToDecimalsByRef(priceInDecimals, token1Decimal-token0Decimal)
	return priceInDecimals
}

func ToDecimalsByRef(reused *big.Float, decimals int) {
	reused.Quo(reused, big.NewFloat(math.Pow10(decimals)))
}

func ToDecimals(weiBalance *big.Float, decimals int) *big.Float {
	ethBalance := new(big.Float)
	ethBalance = ethBalance.Quo(weiBalance, big.NewFloat(math.Pow10(decimals)))
	return ethBalance
}
