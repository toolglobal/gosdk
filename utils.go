package gosdk

import (
	"encoding/hex"
	"strings"
)

func HexToBytes(str string) []byte {
	str = strings.TrimPrefix(str, "0x")
	b, _ := hex.DecodeString(str)
	return b
}

func BytesToHex(bz []byte) string {
	return hex.EncodeToString(bz)
}
