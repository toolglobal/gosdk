package gosdk

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/toolglobal/gosdk/types"
	"math/big"
	"strings"
)

// ERC20Pay
// ERC20代币转账，返回交易hash、使用的gas数量、错误信息
func (cli *APIClient) ERC20Pay(mode string, senderPubKey, senderPriKey, contractAddr, receiver,
	value string, gasLimit uint64, gasPrice string, memo string) (string, uint64, error) {
	if !types.ValidAddress(receiver) {
		return "", 0, errors.New("bad receiver")
	}

	if !types.ValidAddress(contractAddr) {
		return "", 0, errors.New("bad contractAddr")
	}

	if len(memo) > 256 {
		return "", 0, errors.New("bad memo")
	}

	abiIns, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return "", 0, err
	}

	_to := common.HexToAddress(receiver)
	_value, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return "", 0, errors.New("ignore value")
	}
	bin, err := abiIns.Pack("transfer", _to, _value)
	if err != nil {
		return "", 0, err
	}

	tx, _, gasUsed, err := cli.InvokeTx(mode, senderPubKey, senderPriKey, contractAddr,
		"0", BytesToHex(bin), gasLimit, gasPrice, memo)
	return tx, gasUsed, err
}

// ERC20BalanceOf
// 查询ERC20代币账户余额
// 此方法不消耗gas，不上链
func (cli *APIClient) ERC20BalanceOf(senderPubKey, senderPriKey, contractAddr, address string) (*big.Int, error) {
	abiIns, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return nil, err
	}

	bin, err := abiIns.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	data, _, err := cli.QueryTx(senderPubKey, senderPriKey, contractAddr,
		"0", BytesToHex(bin), 100000, "1")
	if err != nil {
		return nil, err
	}

	results, err := abiIns.Unpack("balanceOf", HexToBytes(data))
	if err != nil {
		return nil, err
	}

	return results[0].(*big.Int), err
}
