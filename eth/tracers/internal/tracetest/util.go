package tracetest

import (
	"math/big"
	"strings"
	"unicode"

	"github.com/otherview/go-ethereum/common"
	"github.com/otherview/go-ethereum/common/math"
	"github.com/otherview/go-ethereum/consensus/misc/eip4844"
	"github.com/otherview/go-ethereum/core"
	"github.com/otherview/go-ethereum/core/vm"

	// Force-load native and js packages, to trigger registration
	_ "github.com/otherview/go-ethereum/eth/tracers/js"
	_ "github.com/otherview/go-ethereum/eth/tracers/native"
)

// camel converts a snake cased input string into a camel cased output.
func camel(str string) string {
	pieces := strings.Split(str, "_")
	for i := 1; i < len(pieces); i++ {
		pieces[i] = string(unicode.ToUpper(rune(pieces[i][0]))) + pieces[i][1:]
	}
	return strings.Join(pieces, "")
}

type callContext struct {
	Number     math.HexOrDecimal64   `json:"number"`
	Difficulty *math.HexOrDecimal256 `json:"difficulty"`
	Time       math.HexOrDecimal64   `json:"timestamp"`
	GasLimit   math.HexOrDecimal64   `json:"gasLimit"`
	Miner      common.Address        `json:"miner"`
	BaseFee    *math.HexOrDecimal256 `json:"baseFeePerGas"`
}

func (c *callContext) toBlockContext(genesis *core.Genesis) vm.BlockContext {
	context := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    c.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(c.Number)),
		Time:        uint64(c.Time),
		Difficulty:  (*big.Int)(c.Difficulty),
		GasLimit:    uint64(c.GasLimit),
	}
	if genesis.Config.IsLondon(context.BlockNumber) {
		context.BaseFee = (*big.Int)(c.BaseFee)
	}
	if genesis.ExcessBlobGas != nil && genesis.BlobGasUsed != nil {
		excessBlobGas := eip4844.CalcExcessBlobGas(*genesis.ExcessBlobGas, *genesis.BlobGasUsed)
		context.BlobBaseFee = eip4844.CalcBlobFee(excessBlobGas)
	}
	return context
}
