package gosdk

import (
	"encoding/hex"
	"github.com/axengine/gosdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

// GenKey 本地生成mondo链格式的账户
// 返回：公钥，地址，私钥，mondo采用压缩公钥
func GenKey() (string, string, string, error) {
	privkey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}

	buff := make([]byte, 32)
	copy(buff[32-len(privkey.D.Bytes()):], privkey.D.Bytes())
	return ethcmn.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey)),
		crypto.PubkeyToAddress(privkey.PublicKey).String(),
		ethcmn.Bytes2Hex(buff),
		nil
}

// ValidAddress
// 校验地址是否合法
func ValidAddress(address string) bool {
	return types.ValidAddress(address)
}

// ValidPublicKey
// 校验是否是合法公钥
func ValidPublicKey(publicKey string) bool {
	return types.ValidPublicKey(publicKey)
}

func HexToBytes(str string) []byte {
	str = strings.TrimPrefix(str, "0x")
	b, _ := hex.DecodeString(str)
	return b
}

func BytesToHex(bz []byte) string {
	return hex.EncodeToString(bz)
}
