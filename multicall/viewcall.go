package multicall

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"regexp"
	"strings"
)

type ViewCall struct {
	Key       string
	Target    string
	Method    string
	Arguments []interface{}
}

type ViewCalls []ViewCall

var insideParens = regexp.MustCompile("\\(.*?\\)")

func (vc ViewCall) argumentTypes() []string {
	rawArgs := insideParens.FindAllString(vc.Method, -1)[0]
	rawArgs = strings.Replace(rawArgs, "(", "", -1)
	rawArgs = strings.Replace(rawArgs, ")", "", -1)
	if rawArgs == "" {
		return []string{}
	}
	args := strings.Split(rawArgs, ",")
	for index, arg := range args {
		args[index] = strings.Trim(arg, " ")
	}
	return args
}

func (vc ViewCall) returnTypes() []string {
	rawArgs := insideParens.FindAllString(vc.Method, -1)[1]
	rawArgs = strings.Replace(rawArgs, "(", "", -1)
	rawArgs = strings.Replace(rawArgs, ")", "", -1)
	args := strings.Split(rawArgs, ",")
	for index, arg := range args {
		args[index] = strings.Trim(arg, " ")
	}
	return args
}

func (call ViewCall) callData() ([]byte, error) {
	argsSuffix, err := call.argsCallData()
	if err != nil {
		return nil, err
	}
	methodPrefix, err := call.methodCallData()
	if err != nil {
		return nil, err
	}


	payload := make([]byte, 0)
	payload = append(payload, methodPrefix...)
	payload = append(payload, argsSuffix...)

	return payload, nil
}

func (call ViewCall) methodCallData() ([]byte, error) {
	methodParts := strings.Split(call.Method, ")(")
	var method string
	if len(methodParts) > 1 {
		method = fmt.Sprintf("%s)", methodParts[0])
	} else {
		method = methodParts[0]
	}
	hash := crypto.Keccak256([]byte(method))
	return hash[0:4], nil
}

func (call ViewCall) argsCallData() ([]byte, error) {
	argTypes := call.argumentTypes()
	argValues := make([]interface{}, len(call.Arguments))
	if len(argTypes) != len(call.Arguments) {
		return nil, fmt.Errorf("number of argument types doesn't match with number of arguments for %s with method %s", call.Key, call.Method)
	}
	args := make(abi.Arguments, 0, 0)
	for index, argTypeStr := range argTypes {
		argType, err := abi.NewType(argTypeStr, "", nil)
		if err != nil {
			return nil, err
		}

		args = append(args, abi.Argument{Type: argType})

		if argTypeStr == "address" {
			argValues[index], err = toByteArray(call.Arguments[index].(string))
			if err != nil {
				return nil, err
			}
		} else {
			argValues[index] = call.Arguments[index]
		}
	}
	return args.Pack(argValues...)
}

func (call ViewCall) decode(raw []byte) ([]interface{}, error) {
	retTypes := call.returnTypes()
	args := make(abi.Arguments, 0, 0)
	for index, retTypeStr := range retTypes {
		retType, err := abi.NewType(retTypeStr, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, abi.Argument{Name: fmt.Sprintf("ret%d", index), Type: retType})
	}
	decoded := make(map[string]interface{})
	err := args.UnpackIntoMap(decoded, raw)
	if err != nil {
		return nil, err
	}
	returns := make([]interface{}, len(retTypes))
	for index, _ := range retTypes {
		key := fmt.Sprintf("ret%d", index)
		returns[index] = decoded[key]
	}
	return returns, nil
}

type callArgs struct {
	Target [20]byte
	CallData []byte
}

func (calls ViewCalls) callData() ([]byte, error) {
	payloadArgs := make([]callArgs, 0, len(calls))
	for _, call := range calls {
		callData, err := call.callData()
		if err != nil {
			return nil, err
		}
		targetBytes, err := toByteArray(call.Target)
		if err != nil {
			return nil, err
		}
		payloadArgs = append(payloadArgs, callArgs{targetBytes, callData})
	}

	tupleArray, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Type: "address", Name: "Target"},
		{Type: "bytes", Name: "CallData"},
	})
	if err != nil {
		return nil, err
	}
	boolean, err := abi.NewType("bool", "", nil)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{
		{Type: tupleArray, Name: "calls"},
		{Type: boolean, Name: "strict" },
	}
	return args.Pack(payloadArgs, false)
}

type retType struct {
	Success bool
	Data []byte
}

type wrapperRet struct {
	BlockNumber *big.Int
	Returns []retType
}

func (calls ViewCalls) decode(raw string) (*Result, error) {
	rawBytes, err := hex.DecodeString(strings.Replace(raw, "0x", "", -1))
	if err != nil {
		return nil, err
	}
	result := &Result{}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}
	returnType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "Success", Type: "bool"},
		{Name: "Data", Type: "bytes"},
	})
	if err != nil {
		return nil, err
	}
	wrapperArgs := abi.Arguments{
		{
			Name: "BlockNumber",
			Type: uint256Type,
		},
		{
			Name: "Returns",
			Type: returnType,
		},
	}
	decoded := wrapperRet{}
	err = wrapperArgs.Unpack(&decoded, rawBytes)
	if err != nil {
		return nil, err
	}

	result.BlockNumber = decoded.BlockNumber.Uint64()
	result.Calls = make(map[string]CallResult)
	for index, call := range calls {
		returnValues, err := call.decode(decoded.Returns[index].Data)
		if err != nil {
			return nil, err
		}
		callResult := CallResult{
			Success: decoded.Returns[index].Success,
			ReturnValues: returnValues,
		}
		result.Calls[call.Key] = callResult
	}

	return result, nil
}

func toByteArray(address string) ([20]byte, error) {
	var addressBytes [20]byte;
	address = strings.Replace(address, "0x", "", -1)
	addressBytesSlice, err := hex.DecodeString(address)
	if err != nil {
		return addressBytes, err
	}

	copy(addressBytes[:], addressBytesSlice[:])
	return addressBytes, nil
}

