package main

import (
	"math/rand"
	"log"
	"os"
	"runtime/pprof"
	"github.com/fmstephe/matching_engine/matcher"
	"github.com/fmstephe/matching_engine/trade"
	"time"
	"flag"
)

const (
	stockId = uint32(1)
)

var (
	profile = flag.String("profile", "", "Write out a profile of this application, 'cpu' and 'mem' supported")
)

func main() {
	flag.Parse()
	orderNum := 10 * 1000 * 1000
	sells := mkSells(orderNum, 1000, 1500)
	buys := mkBuys(orderNum, 1000, 1500)
	m := matcher.New(stockId)
	startProfile()
	defer endProfile()
	start := time.Now().UnixNano()
	for i := 0; i < orderNum; i++ {
		m.AddBuy(buys[i])
		m.AddSell(sells[i])
	}
	total := time.Now().UnixNano() - start
	println(total)
}

func startProfile() {
	if *profile == "cpu" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func endProfile() {
	if *profile == "cpu" {
		pprof.StopCPUProfile()
	}
	if *profile == "mem" {
		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
	}
}

func valRangeFlat(n int, low, high int64) []int64 {
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Int63n(high-low) + low
	}
	return vals
}

func valRangePyramid(n int, low, high int64) []int64 {
	seq := (high - low) / 4
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		val := rand.Int63n(seq) + rand.Int63n(seq) + rand.Int63n(seq) + rand.Int63n(seq)
		vals[i] = val + low
	}
	return vals
}

func mkBuys(n int, low, high int64) []*trade.Order {
	return mkOrders(n, low, high, trade.BUY)
}

func mkSells(n int, low, high int64) []*trade.Order {
	return mkOrders(n, low, high, trade.SELL)
}

func mkOrders(n int, low, high int64, buySell trade.TradeType) []*trade.Order {
	prices := valRange(n, low, high)
	orders := make([]*trade.Order, n)
	for i, price := range prices {
		responseFunc := func(response *trade.Response) {
			// Do Nothing
		}
		costData := trade.CostData{Price: price, Amount: 1}
		tradeData := trade.TradeData{TraderId: uint32(i), TradeId: uint32(i), StockId: stockId}
		orders[i] = trade.NewOrder(costData, tradeData, responseFunc, buySell)
	}
	return orders
}

func valRange(n int, low, high int64) []int64 {
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Int63n(high - low) + low
	}
	return vals
}
