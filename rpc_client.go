package gosdk

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/wolot/gosdk/types"
	"time"
)

type RPCClient struct {
	rpcClient *http.HTTP
}

// NewRPCClient
// new tendermint rpc client,remote is tendermint rpc address,like 127.0.0.1:26657
func NewRPCClient(remote string) *RPCClient {
	cli, err := http.New("http://"+remote, "/websocket")
	if err != nil {
		panic(err)
	}
	return &RPCClient{
		rpcClient: cli,
	}
}

// GetValidator
// 获取validator详细信息
func (cli *RPCClient) GetValidator(ctx context.Context, address string) (*types.BeanValidator, error) {
	addressBytes, _ := hex.DecodeString(address)
	resp, err := cli.rpcClient.ABCIQuery(ctx, "/v3/node/account", addressBytes)
	if err != nil {
		return nil, err
	}
	var result types.Result
	err = rlp.DecodeBytes(resp.Response.Value, &result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, errors.New(result.Log)
	}
	fmt.Println(string(result.Data))
	var val types.BeanValidator
	err = json.Unmarshal(result.Data, &val)
	return &val, err
}

// SetRunfor
// 设置节点竞选超级节点标识flag 1：参与竞选 2：不参与竞选
func (cli *RPCClient) SetRunfor(ctx context.Context, genesisPublicKey, genesisPrivateKey string, nodeAddress string, flag uint64) (string, error) {
	var hash string
	tx := types.NewMgrTx()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, err := base64.StdEncoding.DecodeString(genesisPublicKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPublicKey")
	}
	copy(tx.Sender[:], b)

	val, err := cli.GetValidator(ctx, tx.Sender.Address().String())
	if err != nil {
		return hash, errors.Wrap(err, "GetValidator")
	}
	tx.Nonce = val.Nonce

	nodeAddressBytes, err := base64.StdEncoding.DecodeString(nodeAddress)
	if err != nil {
		return hash, errors.Wrap(err, "nodeAddress")
	}
	copy(tx.Receiver[:], nodeAddressBytes)
	tx.OpType = types.OpType_SetRunFor
	tx.OpValue = flag
	b, err = base64.StdEncoding.DecodeString(genesisPrivateKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPrivateKey")
	}
	tx.Signature = tx.Sign(b)

	hash = tx.Hash().Hex()

	result, err := cli.rpcClient.BroadcastTxCommit(ctx, append(types.TxTagAppMgr, tx.ToBytes()...))
	if err != nil {
		return hash, err
	}
	if result.CheckTx.Code != 0 {
		return hash, errors.New(result.CheckTx.Log)
	}

	if result.DeliverTx.Code != 0 {
		return hash, errors.New(result.DeliverTx.Log)
	}
	return hash, nil
}

// SetPower 系统维护使用
// 设置节点的共识投票权power，如果power=0，节点将变成同步节点
func (cli *RPCClient) SetPower(ctx context.Context, genesisPublicKey, genesisPrivateKey string, nodeAddress string, power uint64) (string, error) {
	var hash string
	tx := types.NewMgrTx()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, err := base64.StdEncoding.DecodeString(genesisPublicKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPublicKey")
	}
	copy(tx.Sender[:], b)

	val, err := cli.GetValidator(ctx, tx.Sender.Address().String())
	if err != nil {
		return hash, errors.Wrap(err, "GetValidator")
	}
	tx.Nonce = val.Nonce

	nodeAddressBytes, err := base64.StdEncoding.DecodeString(nodeAddress)
	if err != nil {
		return hash, errors.Wrap(err, "nodeAddress")
	}
	copy(tx.Receiver[:], nodeAddressBytes)
	tx.OpType = types.OpType_SetPower
	tx.OpValue = power
	b, err = base64.StdEncoding.DecodeString(genesisPrivateKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPrivateKey")
	}
	tx.Signature = tx.Sign(b)

	hash = tx.Hash().Hex()

	result, err := cli.rpcClient.BroadcastTxCommit(ctx, append(types.TxTagAppMgr, tx.ToBytes()...))
	if err != nil {
		return hash, err
	}
	if result.CheckTx.Code != 0 {
		return hash, errors.New(result.CheckTx.Log)
	}

	if result.DeliverTx.Code != 0 {
		return hash, errors.New(result.DeliverTx.Log)
	}
	return hash, nil
}

// SetGenesis 系统维护使用
// 转移节点的创世权，调用此接口需慎重，一旦忘记接收者私钥将会丢失创世权
func (cli *RPCClient) SetGenesis(ctx context.Context, genesisPublicKey, genesisPrivateKey string, nodeAddress string) (string, error) {
	var hash string
	tx := types.NewMgrTx()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, err := base64.StdEncoding.DecodeString(genesisPublicKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPublicKey")
	}
	copy(tx.Sender[:], b)

	val, err := cli.GetValidator(ctx, tx.Sender.Address().String())
	if err != nil {
		return hash, errors.Wrap(err, "GetValidator")
	}
	tx.Nonce = val.Nonce

	nodeAddressBytes, err := base64.StdEncoding.DecodeString(nodeAddress)
	if err != nil {
		return hash, errors.Wrap(err, "nodeAddress")
	}
	copy(tx.Receiver[:], nodeAddressBytes)
	tx.OpType = types.OpType_SetGenesis
	tx.OpValue = 0
	b, err = base64.StdEncoding.DecodeString(genesisPrivateKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPrivateKey")
	}
	tx.Signature = tx.Sign(b)

	hash = tx.Hash().Hex()

	result, err := cli.rpcClient.BroadcastTxCommit(ctx, append(types.TxTagAppMgr, tx.ToBytes()...))
	if err != nil {
		return hash, err
	}
	if result.CheckTx.Code != 0 {
		return hash, errors.New(result.CheckTx.Log)
	}

	if result.DeliverTx.Code != 0 {
		return hash, errors.New(result.DeliverTx.Log)
	}
	return hash, nil
}

// SetParam 系统维护使用
func (cli *RPCClient) SetParam(ctx context.Context, genesisPublicKey, genesisPrivateKey string, key string, value string) (string, error) {
	var hash string
	tx := types.NewTxParams()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, err := base64.StdEncoding.DecodeString(genesisPublicKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPublicKey")
	}
	copy(tx.Sender[:], b)

	val, err := cli.GetValidator(ctx, tx.Sender.Address().String())
	if err != nil {
		return hash, errors.Wrap(err, "GetValidator")
	}
	tx.Nonce = val.Nonce

	tx.Key = []byte(key)
	tx.Value = []byte(value)
	b, err = base64.StdEncoding.DecodeString(genesisPrivateKey)
	if err != nil {
		return hash, errors.Wrap(err, "genesisPrivateKey")
	}
	tx.Signature = tx.Sign(b)

	hash = tx.Hash().Hex()

	result, err := cli.rpcClient.BroadcastTxCommit(ctx, append(types.TxTagAppParams, tx.ToBytes()...))
	if err != nil {
		return hash, err
	}
	if result.CheckTx.Code != 0 {
		return hash, errors.New(result.CheckTx.Log)
	}

	if result.DeliverTx.Code != 0 {
		return hash, errors.New(result.DeliverTx.Log)
	}
	return hash, nil
}
