package gosdk

import (
	"context"
	"errors"
	"github.com/axengine/httpc"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/toolglobal/gosdk/types"
	"math/big"
	"strings"
	"time"
)

type APIClient struct {
	addr string
}

// NewAPIClient
// new mondo api client,addr is mondo api http base url,like http://127.0.0.1:8889
func NewAPIClient(addr string) *APIClient {
	return &APIClient{
		addr: addr,
	}
}

const (
	MODE_DEFAULT = "commit"
	MODE_ASYNC   = "async"
	MODE_SYNC    = "sync"
)

func getMode(mode string) int {
	switch mode {
	case MODE_ASYNC:
		return 1
	case MODE_SYNC:
		return 2
	}
	return 0
}

// GetBalance 获取账户余额
// 输入：账户地址或者公钥
// 输出：如果账户未上链返回余额0，error为nil；如果error不为空表示出错
func (cli *APIClient) GetBalance(address string) (*big.Int, error) {
	balance := new(big.Int)
	if len(address) > 42 {
		pub, err := types.HexToPubkey(address)
		if err != nil {
			return balance, err
		}
		address = pub.ToAddress().Hex()
	} else {
		if !types.ValidAddress(address) {
			return balance, errors.New("invalid address")
		}
	}

	bal, _, err := cli.queryAccount(address)
	if err != nil {
		return balance, err
	}
	balance, _ = balance.SetString(bal, 10)
	return balance, nil
}

// Exist 用户是否在链上存在
// 输入：地址或者公钥
func (cli *APIClient) Exist(address string) (bool, error) {
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
	}

	if len(address) > 42 {
		pub, err := types.HexToPubkey(address)
		if err != nil {
			return false, err
		}
		address = pub.ToAddress().Hex()
	}

	err := httpc.New(cli.addr).Path("/v2/accounts/"+address).Get(&resp, httpc.TypeApplicationJson)
	if err != nil {
		return false, err
	}
	if !resp.IsSuccess {
		return false, nil
	}
	return true, nil
}

// GetNonce
// 获取账户nonce
func (cli *APIClient) GetNonce(address string) uint64 {
	if len(address) > 42 {
		pub, err := types.HexToPubkey(address)
		if err != nil {
			return 0
		}
		address = pub.ToAddress().Hex()
	}
	nonce, _ := cli.getNonce(address)
	return nonce
}

func (cli *APIClient) getNonce(address string) (uint64, error) {
	_, nonce, err := cli.queryAccount(address)
	if err != nil {
		return nonce, err
	}
	return nonce, nil
}

// GetAccount
// 获取账户信息，返回余额、nonce，如果账户未上链，返回余额为0，nonce为0
func (cli *APIClient) GetAccount(address string) (string, uint64, error) {
	return cli.queryAccount(address)
}

func (cli *APIClient) queryAccount(address string) (string, uint64, error) {
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
			Nonce   uint64 `json:"nonce"`
		} `json:"result"`
	}
	err := httpc.New(cli.addr).Path("/v2/accounts/"+address).Get(&resp, httpc.TypeApplicationJson)
	if err != nil {
		return "", 0, err
	}
	if !resp.IsSuccess {
		// 未找到用户
		if strings.Contains(resp.Message, "code 11") {
			return "0", 0, nil
		}
		return "", 0, errors.New(resp.Message)
	}
	return resp.Result.Balance, resp.Result.Nonce, nil
}

// GetContract 获取合约
// 输入：合约地址
// 输出：余额、nonce、合约字节码、是否已自杀、错误信息；如果合约不存在，返回错误;如果是普通地址，正确返回；
func (cli *APIClient) GetContract(address string) (*big.Int, uint64, string, bool, error) {
	if !types.ValidAddress(address) {
		return new(big.Int), 0, "", false, errors.New("invalid address")
	}
	return cli.queryContractAccount(address)
}

func (cli *APIClient) queryContractAccount(address string) (*big.Int, uint64, string, bool, error) {
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Address  string `json:"address"`
			Balance  string `json:"balance"`
			Nonce    uint64 `json:"nonce"`
			Code     string `json:"code"`
			Suicided bool   `json:"suicided"`
		} `json:"result"`
	}
	err := httpc.New(cli.addr).Path("/v2/contract/accounts/"+address).Get(&resp, httpc.TypeApplicationJson)
	if err != nil {
		return new(big.Int), 0, "", false, err
	}
	if !resp.IsSuccess {
		return new(big.Int), 0, "", false, errors.New(resp.Message)
	}
	bal, _ := new(big.Int).SetString(resp.Result.Balance, 10)
	return bal, resp.Result.Nonce, resp.Result.Code, resp.Result.Suicided, nil
}

// Payment 发送签名交易
// mode:交易模式 默认为commit模式：交易上链后返回 async:异步模式，节点成功接收交易后返回 sync:同步模式，节点对基本参数校验完成后返回
// senderPubKey:转出方公钥 senderPriKey:转出方私钥 receiver：接收方公钥或地址 value：转账金额（实际金额*1e8）
// gasLimit：gas限额，普通转账21000  gasPrice：建议值1
// 返回交易hash和错误信息
func (cli *APIClient) Payment(mode string, senderPubKey, senderPriKey, receiver, value string, gasLimit uint64, gasPrice string, memo string) (string, error) {
	var (
		hash     string
		err      error
		tx       = types.NewTxEvm()
		signedTx SignedEvmTx
	)

	if len(memo) > 256 {
		return hash, errors.New("bad memo")
	}
	if !types.ValidAddress(receiver) && !types.ValidPublicKey(receiver) {
		return hash, errors.New("bad receiver")
	}
	signedTx.Mode = getMode(mode)
	signedTx.CreatedAt = uint64(time.Now().UnixNano())
	signedTx.GasLimit = gasLimit
	signedTx.GasPrice = gasPrice

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return hash, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	nonce, err := cli.getNonce(tx.Sender.ToAddress().Hex())
	if err != nil {
		return hash, err
	}

	signedTx.Nonce = nonce
	signedTx.Sender = senderPubKey
	signedTx.Body.To = receiver
	signedTx.Body.Value = value
	signedTx.Body.Load = ""
	signedTx.Body.Memo = memo

	tx.CreatedAt = signedTx.CreatedAt
	tx.GasLimit = signedTx.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(signedTx.GasPrice, 10)
	tx.Nonce = signedTx.Nonce

	if len(signedTx.Body.To) > 42 {
		pub, err := types.HexToPubkey(signedTx.Body.To)
		if err != nil {
			return hash, err
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(signedTx.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(signedTx.Body.Value, 10)
	tx.Body.Load, _ = hexutil.Decode(signedTx.Body.Load)
	tx.Body.Memo = []byte(signedTx.Body.Memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return hash, err
	}
	tx.Signature = b
	signedTx.Signature = hexutil.Encode(b)

	// 先计算hash
	hash = tx.Hash().Hex()
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed int         `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/transactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return hash, err
	}
	if !resp.IsSuccess {
		return hash, errors.New(resp.Message)
	}
	return hash, nil
}

// Payments 发送批量交易 1vN
// gasLimit:len(pays)*21000
func (cli *APIClient) Payments(mode string, senderPubKey, senderPriKey string, pays []Operation, gasLimit uint64, gasPrice string, memo string) (string, error) {
	var (
		hash     string
		err      error
		tx       = types.NewTxBatch()
		signedTx SignedBatchTx
	)

	if len(memo) > 256 {
		return hash, errors.New("bad memo")
	}
	signedTx.Mode = getMode(mode)
	signedTx.CreatedAt = uint64(time.Now().UnixNano())
	signedTx.GasLimit = gasLimit
	signedTx.GasPrice = gasPrice

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return hash, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	nonce, err := cli.getNonce(tx.Sender.ToAddress().Hex())
	if err != nil {
		return hash, err
	}

	signedTx.Nonce = nonce
	signedTx.Sender = senderPubKey
	for _, v := range pays {
		signedTx.Ops = append(signedTx.Ops, Operation{
			To:    v.To,
			Value: v.Value,
		})
	}
	signedTx.Memo = memo

	tx.CreatedAt = signedTx.CreatedAt
	tx.GasLimit = signedTx.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(signedTx.GasPrice, 10)
	tx.Nonce = signedTx.Nonce

	for _, v := range signedTx.Ops {
		var op types.TxOp
		if len(v.To) > 42 {
			pub, err := types.HexToPubkey(v.To)
			if err != nil {
				return hash, err
			}
			op.To.SetBytes(pub.Bytes())
		} else {
			if !types.ValidAddress(v.To) {
				return hash, errors.New("bad receiver")
			}
			op.To.SetBytes(ethcmn.HexToAddress(v.To).Bytes())
		}
		op.Value, _ = new(big.Int).SetString(v.Value, 10)

		tx.Ops = append(tx.Ops, op)
	}

	tx.Memo = []byte(signedTx.Memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return hash, err
	}
	tx.Signature = b
	signedTx.Signature = hexutil.Encode(b)

	if int64(gasLimit) < tx.GasWanted() {
		return hash, errors.New("gasLimit too small")
	}

	// 先计算hash
	hash = tx.Hash().Hex()

	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed int         `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/transactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return hash, err
	}
	if !resp.IsSuccess {
		return hash, errors.New(resp.Message)
	}
	return hash, nil
}

// DeployTx 部署合约
// 返回交易hash，合约地址，gas消耗,错误信息
// 大多数错误evm错误都是gas不足导致的，排查问题时首先提高gasLimit
func (cli *APIClient) DeployTx(mode string, senderPubKey, senderPriKey,
	value string, load string, gasLimit uint64, gasPrice string, memo string) (string, string, uint64, error) {
	var (
		hash     string
		err      error
		tx       = types.NewTxEvm()
		signedTx SignedEvmTx
	)

	if len(memo) > 256 {
		return hash, "", 0, errors.New("bad memo")
	}
	signedTx.Mode = getMode(mode)
	signedTx.CreatedAt = uint64(time.Now().UnixNano())
	signedTx.GasLimit = gasLimit
	signedTx.GasPrice = gasPrice

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return hash, "", 0, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	nonce, err := cli.getNonce(tx.Sender.ToAddress().Hex())
	if err != nil {
		return hash, "", 0, err
	}

	contractAddr := crypto.CreateAddress(ethcmn.HexToAddress(tx.Sender.ToAddress().Hex()), nonce).Hex()

	signedTx.Nonce = nonce
	signedTx.Sender = senderPubKey
	signedTx.Body.To = ""
	signedTx.Body.Value = value
	signedTx.Body.Load = load
	signedTx.Body.Memo = memo

	tx.CreatedAt = signedTx.CreatedAt
	tx.GasLimit = signedTx.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(signedTx.GasPrice, 10)
	tx.Nonce = signedTx.Nonce

	if len(signedTx.Body.To) > 42 {
		pub, err := types.HexToPubkey(signedTx.Body.To)
		if err != nil {
			return hash, "", 0, err
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(signedTx.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(signedTx.Body.Value, 10)
	tx.Body.Load = HexToBytes(signedTx.Body.Load)
	tx.Body.Memo = []byte(signedTx.Body.Memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return hash, contractAddr, 0, err
	}
	tx.Signature = b
	signedTx.Signature = hexutil.Encode(b)

	// 先计算hash
	hash = tx.Hash().Hex()
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed uint64      `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/transactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return hash, contractAddr, 0, err
	}
	if !resp.IsSuccess {
		return hash, contractAddr, resp.Result.GasUsed, errors.New(resp.Message)
	}

	return hash, contractAddr, resp.Result.GasUsed, nil
}

// InvokeTx 调用合约
// 返回 交易hash，合约执行返回字节码，gas消耗，错误信息
func (cli *APIClient) InvokeTx(mode string, senderPubKey, senderPriKey, receiver, value string, load string,
	gasLimit uint64, gasPrice string, memo string) (string, string, uint64, error) {
	var (
		hash     string
		err      error
		tx       = types.NewTxEvm()
		signedTx SignedEvmTx
	)

	if len(memo) > 256 {
		return hash, "", 0, errors.New("bad memo")
	}
	signedTx.Mode = getMode(mode)
	signedTx.CreatedAt = uint64(time.Now().UnixNano())
	signedTx.GasLimit = gasLimit
	signedTx.GasPrice = gasPrice

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return hash, "", 0, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	nonce, err := cli.getNonce(tx.Sender.ToAddress().Hex())
	if err != nil {
		return hash, "", 0, err
	}

	signedTx.Nonce = nonce
	signedTx.Sender = senderPubKey
	signedTx.Body.To = receiver
	signedTx.Body.Value = value
	signedTx.Body.Load = load
	signedTx.Body.Memo = memo

	tx.CreatedAt = signedTx.CreatedAt
	tx.GasLimit = signedTx.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(signedTx.GasPrice, 10)
	tx.Nonce = signedTx.Nonce

	if len(signedTx.Body.To) > 42 {
		pub, err := types.HexToPubkey(signedTx.Body.To)
		if err != nil {
			return hash, "", 0, err
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(signedTx.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(signedTx.Body.Value, 10)
	tx.Body.Load = HexToBytes(signedTx.Body.Load)
	tx.Body.Memo = []byte(signedTx.Body.Memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return hash, "", 0, err
	}
	tx.Signature = b
	signedTx.Signature = hexutil.Encode(b)

	// 先计算hash
	hash = tx.Hash().Hex()
	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed uint64      `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/transactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return hash, "", 0, err
	}

	if !resp.IsSuccess {
		return hash, "", resp.Result.GasUsed, errors.New(resp.Message)
	}

	return hash, resp.Result.Ret, resp.Result.GasUsed, nil
}

// QueryTx 查询合约
// 用于估算gas消耗或者调用合约的只读方法查数据
// 只在本节点evm内查询，不形成交易，不上链；在不同的节点查询结果可能不一致（节点高度不一致）；
// 返回 合约执行返回字节码,gas消耗，错误信息
func (cli *APIClient) QueryTx(senderPubKey, senderPriKey, receiver, value string, load string, gasLimit uint64, gasPrice string) (string, uint64, error) {
	var (
		err      error
		tx       = types.NewTxEvm()
		signedTx SignedEvmTx
	)

	signedTx.CreatedAt = uint64(time.Now().UnixNano())
	signedTx.GasLimit = gasLimit
	signedTx.GasPrice = gasPrice

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return "", 0, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	signedTx.Nonce = 1 // 不再强制校验nonce
	signedTx.Sender = senderPubKey
	signedTx.Body.To = receiver
	signedTx.Body.Value = value
	signedTx.Body.Load = load
	signedTx.Body.Memo = VERSION

	tx.CreatedAt = signedTx.CreatedAt
	tx.GasLimit = signedTx.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(signedTx.GasPrice, 10)
	tx.Nonce = signedTx.Nonce

	if len(signedTx.Body.To) > 42 {
		pub, err := types.HexToPubkey(signedTx.Body.To)
		if err != nil {
			return "", 0, err
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(signedTx.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(signedTx.Body.Value, 10)
	tx.Body.Load = HexToBytes(signedTx.Body.Load)
	tx.Body.Memo = []byte(signedTx.Body.Memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return "", 0, err
	}
	tx.Signature = b
	signedTx.Signature = hexutil.Encode(b)

	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    uint32 `json:"code"`    // 错误码
			Msg     string `json:"msg"`     // msg
			Ret     string `json:"ret"`     // 返回数据的hex编码
			GasUsed uint64 `json:"gasUsed"` // 消耗的gas
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/query").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return "", 0, err
	}

	if !resp.IsSuccess {
		return "", 0, errors.New(resp.Message)
	}

	if resp.Result.Code != 0 {
		return "", 0, errors.New(resp.Result.Msg)
	}

	return resp.Result.Ret, resp.Result.GasUsed, nil
}

// BuildEvmTx 生成一个EVM交易
func (cli *APIClient) BuildEvmTx(senderPubKey, senderPriKey, receiver, value string, load string,
	gasLimit uint64, gasPrice string, memo string) (*types.TxEvm, error) {

	var tx = types.NewTxEvm()

	if len(memo) > 256 {
		return nil, errors.New("bad memo")
	}

	if len(senderPubKey) > 42 {
		pub, err := types.HexToPubkey(senderPubKey)
		if err != nil {
			return nil, err
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(senderPubKey).Bytes())
	}

	nonce, err := cli.getNonce(tx.Sender.ToAddress().Hex())
	if err != nil {
		return nil, err
	}

	tx.CreatedAt = uint64(time.Now().Unix())
	tx.GasLimit = gasLimit
	tx.GasPrice, _ = new(big.Int).SetString(gasPrice, 10)
	tx.Nonce = nonce

	if len(receiver) > 42 {
		pub, err := types.HexToPubkey(receiver)
		if err != nil {
			return nil, err
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(receiver).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(value, 10)
	tx.Body.Load = HexToBytes(load)
	tx.Body.Memo = []byte(memo)
	b, err := tx.Sign(senderPriKey)
	if err != nil {
		return nil, err
	}
	tx.Signature = b
	return tx, nil
}

// SendEvmTx 通过API发送EVM Tx
func (cli *APIClient) SendEvmTx(mode string, tx *types.TxEvm) error {
	var signedTx SignedEvmTx

	signedTx.Mode = getMode(mode)
	signedTx.CreatedAt = tx.CreatedAt
	signedTx.GasLimit = tx.GasLimit
	signedTx.GasPrice = tx.GasPrice.String()
	signedTx.Nonce = tx.Nonce
	signedTx.Sender = tx.Sender.String()
	signedTx.Body.To = tx.Body.To.ToAddress().String()
	signedTx.Body.Value = tx.Body.Value.String()
	signedTx.Body.Load = BytesToHex(tx.Body.Load)
	signedTx.Body.Memo = string(tx.Body.Memo)
	signedTx.Signature = BytesToHex(tx.Signature)

	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed uint64      `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/transactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return err
	}

	if !resp.IsSuccess {
		return errors.New(resp.Message)
	}

	return nil
}

// CheckTx 检查交易是否已上链,并处理成功
func (cli *APIClient) CheckTx(hash string) (bool, error) {
	if hash == "" {
		return false, errors.New("error hash")
	}
	tx, err := cli.V3GetTransaction(hash)
	if err != nil || tx == nil {
		return false, err
	}
	if tx.Codei != 0 {
		return false, nil
	}
	return true, nil
}

// CheckTxWithCtx 检查交易是否已上链，并处理成功.ctx应待deadline
func (cli *APIClient) CheckTxWithCtx(ctx context.Context, hash string) (bool, error) {
	if hash == "" {
		return false, errors.New("error hash")
	}
	for {
		select {
		case <-ctx.Done():
			return false, errors.New("canceled")
		default:
			tx, err := cli.V3GetTransaction(hash)
			if err != nil || tx == nil {
				continue
			}
			if tx.Codei != 0 {
				return false, nil
			}
			return true, nil
		}
	}
}

// SendMultisigEvmTx 发送多签交易
func (cli *APIClient) SendMultisigEvmTx(mode string, tx *types.MultisigEvmTx) error {
	var signedTx SignedMultisigEvmTx

	signedTx.Mode = getMode(mode)
	signedTx.Deadline = tx.Deadline
	signedTx.GasLimit = tx.GasLimit
	signedTx.GasPrice = tx.GasPrice.String()
	signedTx.Nonce = tx.Nonce
	signedTx.From = tx.From.String()
	signedTx.To = tx.To.String()
	signedTx.Value = tx.Value.String()
	signedTx.Load = BytesToHex(tx.Load)
	signedTx.Memo = string(tx.Memo)

	signedTx.Signature.PubKey.K = int(tx.Signature.PubKey.K)
	for _, v := range tx.Signature.PubKey.PubKeys {
		signedTx.Signature.PubKey.PubKeys = append(signedTx.Signature.PubKey.PubKeys, v.String())
	}

	for _, v := range tx.Signature.MultiSig.Sigs {
		signedTx.Signature.Signatures = append(signedTx.Signature.Signatures, BytesToHex(v))
	}

	var resp struct {
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		Result    struct {
			Code    int         `json:"code"`
			GasUsed uint64      `json:"gasUsed"`
			Logs    interface{} `json:"logs"`
			Ret     string      `json:"ret"`
			Tx      string      `json:"tx"`
		} `json:"result"`
	}
	if err := httpc.New(cli.addr).Path("/v2/contract/multisigTransactions").
		ContentType(httpc.TypeApplicationJson).
		Body(&signedTx).Post(&resp); err != nil {
		return err
	}

	if !resp.IsSuccess {
		return errors.New(resp.Message)
	}

	return nil
}
