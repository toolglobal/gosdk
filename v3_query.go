package gosdk

import (
	"github.com/axengine/httpc"
	"github.com/pkg/errors"
	"strconv"
)

// V3GetLedgers 查询所有账本（区块）
func (cli *APIClient) V3GetLedgers(cursor, limit int64, order string) ([]V3Ledger, error) {
	resp := struct {
		IsSuccess bool       `json:"isSuccess"`
		Result    []V3Ledger `json:"result"`
	}{}
	err := httpc.New(cli.addr).Path("/v3/ledgers").Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New("get ledgers error")
	}

	return resp.Result, nil
}

// V3GetLedger 指定高度查询账本（区块）
func (cli *APIClient) V3GetLedger(height int64) (*V3Ledger, error) {
	resp := struct {
		IsSuccess bool       `json:"isSuccess"`
		Result    []V3Ledger `json:"result"`
	}{}
	err := httpc.New(cli.addr).Path("/v3/ledgers/" + strconv.FormatInt(height, 10)).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New("get ledgers error")
	}

	if len(resp.Result) > 0 {
		return &resp.Result[0], nil
	}

	return nil, nil
}

// V3GetTransactions 查询所有交易
func (cli *APIClient) V3GetTransactions(cursor, limit int64, order string) ([]V3Transaction, error) {
	resp := struct {
		IsSuccess bool            `json:"isSuccess"`
		Message   string          `json:"message"`
		Result    []V3Transaction `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/transactions").Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

// V3GetTransaction 指定hash查交易
func (cli *APIClient) V3GetTransaction(tx string) (*V3Transaction, error) {
	resp := struct {
		IsSuccess bool            `json:"isSuccess"`
		Message   string          `json:"message"`
		Result    []V3Transaction `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/transactions/" + tx).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	if len(resp.Result) == 0 {
		return nil, nil
	}
	return &resp.Result[0], nil
}

// V3GetPayments 查询所有转账记录
// symbol：币种，例如OLO，大写，为空时查询所有币种的转账记录
func (cli *APIClient) V3GetPayments(symbol string, cursor, limit int64, order string) ([]V3Payment, error) {
	resp := struct {
		IsSuccess bool        `json:"isSuccess"`
		Message   string      `json:"message"`
		Result    []V3Payment `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/payments").Query("symbol", symbol).Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

// V1GetTxPayments 指定交易hash查payment
func (cli *APIClient) V3GetTxPayments(tx string, cursor, limit int64, order string) ([]V3Payment, error) {
	resp := struct {
		IsSuccess bool        `json:"isSuccess"`
		Message   string      `json:"message"`
		Result    []V3Payment `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/transactions/"+tx+"/payments").Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

// V3GetAccountPayments 指定账户查payment，账户可以是地址，也可以是公钥
func (cli *APIClient) V3GetAccountPayments(address string, cursor, limit int64, order string) ([]V3Payment, error) {
	resp := struct {
		IsSuccess bool        `json:"isSuccess"`
		Message   string      `json:"message"`
		Result    []V3Payment `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/accounts/"+address+"/payments").Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

// V3GetAccountTransactions 指定账户查payment
func (cli *APIClient) V3GetAccountTransactions(address string, cursor, limit int64, order string) ([]V3Transaction, error) {
	resp := struct {
		IsSuccess bool            `json:"isSuccess"`
		Message   string          `json:"message"`
		Result    []V3Transaction `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/accounts/"+address+"/transactions").Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

// V3GetLedgerTransactions 指定高度查tx
func (cli *APIClient) V3GetLedgerTransactions(height int64, cursor, limit int64, order string) ([]V3Transaction, error) {
	resp := struct {
		IsSuccess bool            `json:"isSuccess"`
		Message   string          `json:"message"`
		Result    []V3Transaction `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/ledgers/"+strconv.FormatUint(uint64(height), 10)+"/transactions").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	if len(resp.Result) == 0 {
		return nil, nil
	}
	return resp.Result, nil
}

// V3GetLedgerPayments 指定高度查账本
func (cli *APIClient) V3GetLedgerPayments(height int64, cursor, limit int64, order string) ([]V3Payment, error) {
	resp := struct {
		IsSuccess bool        `json:"isSuccess"`
		Message   string      `json:"message"`
		Result    []V3Payment `json:"result"`
	}{}

	err := httpc.New(cli.addr).Path("/v3/ledgers/"+strconv.FormatUint(uint64(height), 10)+"/payments").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).Query("order", order).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}
	if len(resp.Result) == 0 {
		return nil, nil
	}
	return resp.Result, nil
}
