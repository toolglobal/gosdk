package types

import (
	"fmt"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func Test_x(t *testing.T) {
	digit := sha3.New256()
	digit.Write([]byte{1, 2, 3, 4})
	hx := digit.Sum(nil)

	pkSet1, sigSet1 := generatePubKeysAndSignatures(3, hx)
	multisigKey := NewPubKeyMultisigThreshold(2, pkSet1)
	multisignature := NewMultisig(3)

	fmt.Println("multisigKey address=", multisigKey.Address().Hex())

	for i := 0; i < 2; i = i + 1 {
		multisignature.AddSignatureFromPubKey(sigSet1[i], pkSet1[i], pkSet1)
	}

	ok := multisigKey.VerifyBytes(hx, multisignature.Marshal())
	if !ok {
		t.Fatal(ok)
	}
}

func generatePubKeysAndSignatures(n int, msg []byte) (pubkeys []PublicKey, signatures [][]byte) {
	pubkeys = make([]PublicKey, n)
	signatures = make([][]byte, n)
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		buff := make([]byte, 32)
		copy(buff[32-len(key.D.Bytes()):], key.D.Bytes())
		privkey := ethcmn.Bytes2Hex(buff)
		pubkey := ethcmn.Bytes2Hex(crypto.CompressPubkey(&key.PublicKey))
		address := crypto.PubkeyToAddress(key.PublicKey).String()

		fmt.Println("privKey:", privkey, " pubKey:", pubkey, " address:", address)
		priv, _ := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
		signature, err := crypto.Sign(msg, priv)
		if err != nil {
			panic(err)
		}

		copy(pubkeys[i][:], ethcmn.Hex2Bytes(pubkey))
		signatures[i] = signature
	}
	return
}
