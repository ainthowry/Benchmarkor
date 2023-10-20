package benchmarkor

import (
	"context"
	"encoding/csv"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"golang.org/x/exp/slog"
)

type Benchmarkor struct {
	maxConcurrentConn uint64
	iterations        uint64
	timeout           time.Duration

	rpcCall RpcCall

	clients      []*RpcClient
	rpcHostnames []string

	wg sync.WaitGroup

	data [][]BenchmarkData
}

type BenchmarkorOpts struct {
	RpcUrls []string

	MaxConcurrentConn uint64
	Iterations        uint64
	Timeout           time.Duration
}

type BenchmarkData struct {
	status    int
	timeTaken uint64
}

func (b *Benchmarkor) SetRpcCall(rpcCall RpcCall) {
	b.rpcCall = rpcCall
}

func (b *Benchmarkor) BenchmarkClient(idx int) {
	poolch := make(chan struct{}, b.maxConcurrentConn)
	var wg sync.WaitGroup
	wg.Add(int(b.iterations))
	parent := context.Background()

	blockNumber := new(big.Int).Set(b.clients[idx].blockNumber)
	for i := uint64(0); i < b.iterations; i++ {
		poolch <- struct{}{}
		ctx, cancel := context.WithTimeout(parent, b.timeout)

		blockNumber.Sub(blockNumber, new(big.Int).SetUint64(i))
		callOpts := &bind.CallOpts{Pending: false, BlockNumber: new(big.Int).Set(blockNumber), Context: ctx}

		iteridx := i
		cancelFunc := cancel
		go func(callOpts *bind.CallOpts, iteridx uint64, cancel context.CancelFunc) {
			defer cancel()
			defer wg.Done()

			status, time, _ := b.clients[idx].rpcCall(b.rpcCall, callOpts)
			b.data[idx][iteridx] = BenchmarkData{status: status, timeTaken: time}
			<-poolch
		}(callOpts, iteridx, cancelFunc)
	}
	wg.Wait()
}

func (b *Benchmarkor) BenchmarkAll() {
	for idx := range b.clients {
		b.wg.Add(1)
		go func(idx int) {
			defer b.wg.Done()
			b.BenchmarkClient(idx)
		}(idx)
	}
	b.wg.Wait()
}

func (b *Benchmarkor) ExportCsv(fileName_optional ...string) error {
	fileName := "results.csv"
	if len(fileName_optional) > 0 {
		fileName = fileName_optional[0]
	}
	file, err := os.Create(fileName)
	if err != nil {
		slog.Warn("Unable to create file", "fileName", fileName, "err", err)
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	for idx, data := range b.data {

		for _, benchmarkResult := range data {
			row := []string{b.rpcHostnames[idx], strconv.Itoa(benchmarkResult.status), strconv.FormatUint(benchmarkResult.timeTaken, 10)}
			if err := w.Write(row); err != nil {
				slog.Warn("Unable to write to file", "err", err)
				return err
			}
		}
	}
	return nil
}

func (b *Benchmarkor) ClearData() {
	for idx := range b.data {
		clear(b.data[idx])
	}
}

func NewBenchmarkor(opts BenchmarkorOpts) (benchmarkor *Benchmarkor, err error) {
	clients := make([]*RpcClient, len(opts.RpcUrls))
	rpcHostnames := make([]string, len(opts.RpcUrls))
	data := make([][]BenchmarkData, len(clients))

	for idx, rpcUrl := range opts.RpcUrls {
		data[idx] = make([]BenchmarkData, opts.Iterations)

		url, err := url.Parse(rpcUrl)
		if err != nil {
			slog.Warn("Not valid URL parsed")
			return nil, err
		}
		rpcHostnames[idx] = url.Hostname()

		clients[idx], err = NewRpcClient(rpcUrl)
		if err != nil {
			slog.Warn("NewRpcClient Failed")
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		blockNumber, err := clients[idx].client.BlockNumber(ctx)
		if err != nil {
			slog.Warn("GetBlock Failed")
			return nil, err
		}
		clients[idx].blockNumber = new(big.Int).SetUint64(blockNumber)
	}

	benchmarkor = &Benchmarkor{
		maxConcurrentConn: opts.MaxConcurrentConn,
		iterations:        opts.Iterations,
		timeout:           opts.Timeout,

		clients:      clients,
		rpcHostnames: rpcHostnames,

		data: data,
	}
	return benchmarkor, nil
}
