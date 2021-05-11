package multicall

import "fmt"

type Option func(*Config)

type Config struct {
	MulticallAddress string
	Gas              string
}

const (
	// MainnetMulticall : Multicall contract address on mainnet
	MainnetAddress = "0x5eb3fa2dfecdde21c950813c665e9364fa609bd2"
	// RopstenMulticall : Multicall contract address on Ropsten
	RopstenAddress = "0xf3ad7e31b052ff96566eedd218a823430e74b406"
)

func ContractAddress(address string) Option {
	return func(c *Config) {
		c.MulticallAddress = address
	}
}

func SetGas(gas uint64) Option {
	return func(c *Config) {
		c.Gas = fmt.Sprintf("0x%x", gas)
	}
}

func SetGasHex(gas string) Option {
	return func(c *Config) {
		c.Gas = gas
	}
}
