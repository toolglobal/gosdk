package gosdk

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/wolot/gosdk/types"
	"math/big"
	"strings"
)

// ERC20Pay
func (api *APISDK) ERC20Pay(mode string, senderPubKey, senderPrivKey, contractAddr, receiver,
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

	tx, _, gasUsed, err := api.InvokeTx(mode, senderPubKey, senderPrivKey, contractAddr,
		"0", BytesToHex(bin), gasLimit, gasPrice, memo)
	return tx, gasUsed, err
}

// ERC20BalanceOf
// 查询ERC20代币账户余额
// 此方法不消耗gas，不上链，但要求sender账户有足够余额
func (api *APISDK) ERC20BalanceOf(senderPubKey, senderPrivKey, contractAddr, address string) (*big.Int, error) {
	balance := new(big.Int)
	abiIns, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return balance, err
	}

	bin, err := abiIns.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		return balance, err
	}

	data, _, err := api.QueryTx(senderPubKey, senderPrivKey, contractAddr,
		"0", BytesToHex(bin), 100000, "1")
	if err != nil {
		return balance, err
	}

	if err := abiIns.Unpack(&balance, "balanceOf", HexToBytes(data)); err != nil {
		return balance, err
	}

	return balance, err
}
