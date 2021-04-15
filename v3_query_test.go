package gosdk

import (
	"fmt"
	"testing"
)

func TestAPISDK_nV3GetLedgers(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetLedgers(0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetLedger(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetLedger(1000000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetTransactions(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetTransactions(0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetTransaction(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetTransaction("0x1a0d24043c31ed25360eaabf019c7c070d0292723d091838ba566c47b2a577f4")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetPayments(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetPayments("ABC", 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetTxPayments(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetTxPayments("0xa8cd9dd3a6c6ef075c23e7a754ed9bb14c5fc0ca215f2b6a636b8e6a262b6f9a", 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetAccountTransactions(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetAccountTransactions("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F", 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetAccountPayments(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetAccountPayments("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F", 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetLedgerPayments(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetLedgerPayments(764114, 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}

func TestAPISDK_V3GetLedgerTransactions(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	datas, err := api.V3GetLedgerTransactions(764114, 0, 10, "ASC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(datas)
}
