package multicall_test

import (
	"encoding/hex"
	"fmt"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/alethio/web3-multicall-go/multicall"
	"time"
)


func ExampleViwCall() {
	eth, err := getETH("https://mainnet.infura.io/")
	vcs := multicall.ViewCalls{
		{
			Key:       "key.1",
			Target:    "0x5eb3fa2dfecdde21c950813c665e9364fa609bd2",
			Method:    "getLastBlockHash()(bytes32)",
			Arguments: []interface{}{},
		},
	}
	mc, _ := multicall.New(eth)
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	blockHash := res.Calls["key.4"].ReturnValues[0].([32]byte)
	fmt.Println(hex.EncodeToString(blockHash[:]))
	fmt.Println(err)

}

func getETH(url string) (ethrpc.ETHInterface, error) {
	provider, err := httprpc.New(url)
	if err != nil {
		return nil, err
	}
	provider.SetHTTPTimeout(5 * time.Second)
	return ethrpc.New(provider)
}

