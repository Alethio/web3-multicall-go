package multicall

import (
	"fmt"
	"math/big"
)

type BigIntJSONString big.Int

func (i BigIntJSONString) MarshalJSON() ([]byte, error) {
	backToInt := big.Int(i)
	return []byte(fmt.Sprintf(`"%s"`, backToInt.String())), nil
}
