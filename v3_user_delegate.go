package gosdk

import (
	"encoding/hex"
	"errors"
	"github.com/axengine/gosdk/types"
	"github.com/axengine/httpc"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

// Delegate
// 用户抵押赎回操作
// opType:1-抵押选举 2-赎回 3-领取收益
// receiverAddress:抵押时必填抵押节点地址
// opValue：操作金额，抵押时为抵押金额，赎回时不填，领取收益时为领取收益金额
// 返回hash和错误信息
func (cli *APIClient) Delegate(publicKey, privateKey, receiverAddress string, opType uint8, opValue string) (string, error) {
	var (
		hash string
		err  error
	)
	var bean DelegateTx
	bean.CreatedAt = uint64(time.Now().Unix())
	bean.Sender = publicKey
	bean.Nonce = cli.GetNonce(publicKey)
	bean.OpType = opType
	bean.OpValue = opValue
	bean.Receiver = receiverAddress
	tx := types.NewUserDelegateTx()
	tx.CreatedAt = bean.CreatedAt
	tx.Sender, _ = types.HexToPubkey(bean.Sender)
	tx.Nonce = bean.Nonce
	tx.OpType = bean.OpType
	if len(bean.Receiver) > 0 {
		tx.Receiver = common.HexToAddress(bean.Receiver).Bytes()
	}
	tx.OpValue, _ = new(big.Int).SetString(bean.OpValue, 10)
	b, err := tx.Sign(privateKey)
	if err != nil {
		return hash, err
	}
	tx.Signature = b
	bean.Signature = hex.EncodeToString(tx.Signature)

	// 先计算hash
	hash = tx.Hash().Hex()
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Tx string `json:"tx"`
		} `json:"result"`
	}
	if err = httpc.New(cli.addr).Path("/v3/delegate").
		ContentType(httpc.TypeApplicationJson).
		Body(&bean).Post(&resp); err != nil {
		return hash, err
	}
	if !resp.IsSuccess {
		return hash, errors.New(resp.Message)
	}
	return hash, nil
}
