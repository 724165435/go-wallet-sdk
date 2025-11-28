package aptos

import (
	"encoding/hex"

	"github.com/724165435/go-wallet-sdk/crypto/ed25519"
	//"github.com/724165435/go-wallet-sdk/wallet"
)

const HexPrefix = "0x"

// type AptosWallet struct {
// 	wallet.WalletBase
// }

func GetRandomPrivateKey() (string, error) {
	p, err := ed25519.GenerateKey()
	if err != nil {
		return "", err
	}
	return HexPrefix + hex.EncodeToString(p), nil
}
