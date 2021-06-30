package gosdk

import (
	"context"
	"fmt"
	"github.com/axengine/gosdk/types"
	"strconv"
	"testing"
)

var (
	NODE_ADDRESS  = "517D364E220D347E6B43190843A8CE8C159D7A7B"
	NODE_PUB_KEY  = "9VcVrno+9fdy6r5fJB0gpJCVRsjy7iA5+UxMX6FH4KM="
	NODE_PRIV_KEY = "nSmQCkjXOdy+B453NpJSKqJkYLaZHT6aZ6Hh7lYdakj1VxWuej7193Lqvl8kHSCkkJVGyPLuIDn5TExfoUfgow=="
)

func TestRPCSDK_GetValidator(t *testing.T) {
	sdk := NewRPCClient("192.168.10.110:26657")
	val, err := sdk.GetValidator(context.Background(), NODE_ADDRESS)
	fmt.Println(val)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}

func TestRPCSDK_SetRunfor(t *testing.T) {
	sdk := NewRPCClient("192.168.10.110:26657")
	hash, err := sdk.SetRunfor(context.Background(), NODE_PUB_KEY, NODE_PRIV_KEY, "lcNeyXdktqjqlsewjz8k7ufmMSLh8Kxqgk0Trd1d8mE=", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}

func TestRPCSDK_SetPower(t *testing.T) {
	sdk := NewRPCClient("192.168.10.110:26657")
	hash, err := sdk.SetPower(context.Background(), NODE_PUB_KEY, NODE_PRIV_KEY, "lcNeyXdktqjqlsewjz8k7ufmMSLh8Kxqgk0Trd1d8mE=", 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}

func TestRPCSDK_SetParam(t *testing.T) {
	sdk := NewRPCClient("192.168.10.110:26657")
	hash, err := sdk.SetParam(context.Background(), NODE_PUB_KEY, NODE_PRIV_KEY, types.KEY_UpgradeHeight, strconv.FormatInt(129860, 10))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}
