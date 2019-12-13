package multicall

import (
	"encoding/hex"
	"github.com/alethio/web3-go/ethrpc"
)

type Multicall struct {
	eth    ethrpc.ETHInterface
	config *Config
}


func New(eth ethrpc.ETHInterface, opts ...Option) (*Multicall, error) {
	config := &Config {
		MulticallAddress: MainnetAddress,
		Gas: "0x400000000",
	}

	for _, opt := range opts {
		opt(config)
	}

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

func (mc *Multicall) Contract() string {
	return mc.config.MulticallAddress
}


