package multicall

import (
	"encoding/hex"
	"github.com/alethio/web3-go/ethrpc"
)

type Multicall struct {
	eth    ethrpc.ETHInterface
	config Config
}

type Config struct {
	MulticallAddress string
	Gas				 string
}

var (
	// MainnetConfig exposes a default config for the mainnet multicall contract
	MainnetConfig = Config{"0x5eb3fa2dfecdde21c950813c665e9364fa609bd2", "0x4000000000000"}
	// RopstenConfig exposes a default config for the ropsten multicall contract
	RopstenConfig = Config{"0xf3ad7e31b052ff96566eedd218a823430e74b406", "0x4000000000000"}
)


func New(eth ethrpc.ETHInterface, config Config) (*Multicall, error) {
	return &Multicall{
		eth: eth,
		config: config,
	}, nil
}

type CallResult struct {
	Success bool
	ReturnValues []interface{}
}

type Result struct {
	BlockNumber uint64
	Calls map[string]CallResult
}

const AggregateMethod = "0x17352e13"

func (mc *Multicall) Call(calls ViewCalls, block string) (*Result, error) {
	payloadArgs, err := calls.callData()
	if err != nil {
		return nil, err
	}
	payload := make(map[string]string)
	payload["to"] = mc.config.MulticallAddress
	payload["data"] = AggregateMethod + hex.EncodeToString(payloadArgs)
	payload["gas"] = mc.config.Gas
	var resultRaw string
	err = mc.eth.MakeRequest(&resultRaw, ethrpc.ETHCall, payload, block)
	if err != nil {
		return nil, err
	}
	return calls.decode(resultRaw)
}


