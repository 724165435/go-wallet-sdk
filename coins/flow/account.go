package flow

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"

	"github.com/724165435/go-wallet-sdk/coins/flow/core"
	"github.com/724165435/go-wallet-sdk/util"
	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
)

func GenerateKeyPair() (privKey, pubKey string) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return hex.EncodeToString(privateKey.D.Bytes()), hex.EncodeToString(pubKeyBytes)
}
// DerivePublicKeyFromPrivate 根据私钥推导公钥（Flow格式）
func DerivePublicKeyFromPrivate(privateKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	// 使用btcec库恢复私钥
	privateKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	privateKeyEcdsa := privateKey.ToECDSA()

	// 提取公钥坐标
	xBytes := privateKeyEcdsa.PublicKey.X.Bytes()
	yBytes := privateKeyEcdsa.PublicKey.Y.Bytes()

	// 填充到32字节
	paddedX := make([]byte, 32)
	paddedY := make([]byte, 32)
	copy(paddedX[32-len(xBytes):], xBytes)
	copy(paddedY[32-len(yBytes):], yBytes)

	// Flow格式：0x04 + X + Y
	publicKeyBytes := append([]byte{0x04}, paddedX...)
	publicKeyBytes = append(publicKeyBytes, paddedY...)

	return hex.EncodeToString(publicKeyBytes), nil
}

// DerivePublicKeyFromPrivateBigInt 使用大整数操作的替代实现
func DerivePublicKeyFromPrivateBigInt(privateKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	// 将私钥字节转换为大整数
	d := new(big.Int).SetBytes(privKeyBytes)

	// 使用P256曲线计算公钥点
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(d.Bytes())

	// 编码公钥点坐标
	xBytes := x.Bytes()
	yBytes := y.Bytes()

	// 填充到32字节
	paddedX := make([]byte, 32)
	paddedY := make([]byte, 32)
	copy(paddedX[32-len(xBytes):], xBytes)
	copy(paddedY[32-len(yBytes):], yBytes)

	// Flow格式：0x04 + X + Y
	publicKeyBytes := append([]byte{0x04}, paddedX...)
	publicKeyBytes = append(publicKeyBytes, paddedY...)

	return hex.EncodeToString(publicKeyBytes), nil
}

// DerivePublicKeyFromPrivateBigIntRaw 使用大整数操作推导公钥（去掉0x04前缀）
func DerivePublicKeyFromPrivateBigIntRaw(privateKeyHex string) (string, error) {
	pubKeyWithPrefix, err := DerivePublicKeyFromPrivateBigInt(privateKeyHex)
	if err != nil {
		return "", err
	}

	return pubKeyWithPrefix[2:], nil
}
func SignTx(signerAddr, privKeyHex string, tx *core.Transaction) error {
	envelopeMessage := tx.EnvelopeMessage()
	transactionDomainTag := new([32]byte)
	copy(transactionDomainTag[:], "FLOW-V0.0-transaction")
	message := append(transactionDomainTag[:], envelopeMessage...)
	hashBytes := hashSha256(message)
	sig, err := signEcdsaP256(hashBytes, privKeyHex)
	if err != nil {
		return err
	}
	signature := core.TransactionSignature{
		Address:     core.HexToAddress(signerAddr),
		SignerIndex: 0,
		KeyIndex:    0,
		Signature:   sig,
	}
	tx.EnvelopeSignatures = []core.TransactionSignature{signature}
	return nil
}

func hashSha256(message []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(message)
	return hasher.Sum(nil)
}

func signEcdsaP256(hash []byte, privateKeyHex string) ([]byte, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	privateKeyEcdsa := privateKey.ToECDSA()
	r, s, err := ecdsa.Sign(rand.Reader, privateKeyEcdsa, hash)
	if err != nil {
		return nil, err
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	Nlen := bitsToBytes((privateKeyEcdsa.PublicKey.Curve.Params().N).BitLen())
	signature := make([]byte, 2*Nlen)
	// pad the signature with zeroes
	copy(signature[Nlen-len(rBytes):], rBytes)
	copy(signature[2*Nlen-len(sBytes):], sBytes)
	return signature, nil
}

func bitsToBytes(bits int) int {
	return (bits + 7) >> 3
}
func ValidateAddress(address string) bool {
	bytes := util.DecodeHexString(address)
	return len(bytes) == 8
}
