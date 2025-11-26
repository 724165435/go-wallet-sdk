package types

import (
	"github.com/724165435/go-wallet-sdk/coins/cosmos/okc/tx/amino"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&StdTx{}, "cosmos-sdk/StdTx", nil)
}

func init() {
	RegisterCodec(amino.GCodec)
}
