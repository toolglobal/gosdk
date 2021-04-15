package gosdk

import (
	"github.com/wolot/gosdk/types"
	"testing"
)

func TestAPISDK_Delegate(t *testing.T) {
	api := NewAPISDK(DEV_API_URL_BASE)
	tx, err := api.Delegate(DEV_USER_PUBKEY,
		DEV_USER_PRIVKEY, "",
		types.USER_OPTYPE_REEDEM, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx)
}
