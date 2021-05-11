package multicall_test

import (
	"encoding/json"
	"fmt"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExampleViwCall(t *testing.T) {
	eth, err := getETH("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")
	vc := multicall.NewViewCall(
		"key.1",
		"0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643",
		"totalReserves()(uint256)",
		[]interface{}{},
	)
	vcs := multicall.ViewCalls{vc}
	mc, _ := multicall.New(eth)
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}

	resJson, _ := json.Marshal(res)
	fmt.Println(string(resJson))
	fmt.Println(res)
	fmt.Println(err)

}

func TestViwCallWithDecodeError(t *testing.T) {
	eth, err := getETH("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")
	vc1 := multicall.NewViewCall(
		"good_call",
		"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
		"symbol()(string)",
		[]interface{}{},
	)
	vc2 := multicall.NewViewCall(
		"broken_call",
		"0xa417221ef64b1549575c977764e651c9fab50141",
		"latestAnswer()(int256)",
		[]interface{}{},
	)

	vcs := multicall.ViewCalls{vc1, vc2}
	mc, _ := multicall.New(eth)
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, true, res.Calls["good_call"].Success)
	assert.Equal(t, "UNI", res.Calls["good_call"].Decoded[0].(string))

	assert.Equal(t, false, res.Calls["broken_call"].Success)
}

func getETH(url string) (ethrpc.ETHInterface, error) {
	provider, err := httprpc.New(url)
	if err != nil {
		return nil, err
	}
	provider.SetHTTPTimeout(5 * time.Second)
	return ethrpc.New(provider)
}
