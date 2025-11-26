package zcash

import "github.com/724165435/go-wallet-sdk/util"

func ValidateAddress(address string) bool {
	_, v, err := util.CheckDecodeDoubleV(address)
	return err == nil && (v == [2]byte{0x1c, 0xb8} || v[1] == 0xbd) && len(address) == 35
}
