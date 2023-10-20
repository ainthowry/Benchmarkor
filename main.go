package main

import (
	"benchmarkor-go/apis"
	"benchmarkor-go/benchmarkor"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

const MAX_CONNECTIONS = 5
const ITERATIONS = 100
const TIMEOUT = time.Second * 10

var RPC_URLS = []string{
	"https://eth-mainnet.g.alchemy.com/v2/xOJLgMkKNGcGSKV9ixPnvHFUU2UUNyf1",
	"https://mainnet.infura.io/v3/56111aa4547d4955967dd6f87c1f9fef",
	"https://fragrant-proud-bush.discover.quiknode.pro/755ee61b256f1af36802a6845fa05cbb64d2c1e4",
}

func init() {
	flag.Parse()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Godotenv", "err", err)
	}

	var rpcs_urls []string
	for _, rpc_url := range strings.Split(os.Getenv("RPCS"), "\n") {
		if rpc_url != "" {
			rpcs_urls = append(rpcs_urls, rpc_url)
		}
	}

	opts := &benchmarkor.BenchmarkorOpts{
		MaxConcurrentConn: MAX_CONNECTIONS,
		Iterations:        ITERATIONS,
		Timeout:           TIMEOUT,

		RpcUrls: rpcs_urls,
	}
	benchmarkor, err := benchmarkor.NewBenchmarkor(*opts)
	if err != nil {
		slog.Error("NewBenchmarkor", "err", err)
	}

	benchmarkor.SetRpcCall(apis.GetBlockByNumber)
	benchmarkor.BenchmarkAll()

	benchmarkor.ExportCsv("results.csv")
	benchmarkor.ClearData()

	os.Exit(0)
}
