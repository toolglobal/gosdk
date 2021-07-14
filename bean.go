package gosdk

type SignedEvmTx struct {
	Mode      int    `json:"mode"`      // 模式:0-default/commit 1-async 2-sync
	CreatedAt uint64 `json:"createdAt"` // 时间戳unixNano
	GasLimit  uint64 `json:"gasLimit"`  //
	GasPrice  string `json:"gasPrice"`  //
	Nonce     uint64 `json:"nonce"`     //
	Sender    string `json:"sender"`    // pubkey
	Body      struct {
		To    string `json:"to"`    // 合约地址
		Value string `json:"value"` //
		Load  string `json:"load"`  // hex编码
		Memo  string `json:"memo"`  // 备注
	} `json:"body"`
	Signature string `json:"signature"` // hex编码
}

// 批量交易
type SignedBatchTx struct {
	Mode      int         `json:"mode"`      // 模式:0-default/commit 1-async 2-sync
	CreatedAt uint64      `json:"createdAt"` // unix nano
	GasLimit  uint64      `json:"gasLimit"`  //
	GasPrice  string      `json:"gasPrice"`  //
	Nonce     uint64      `json:"nonce"`     //
	Sender    string      `json:"sender"`    // address
	Ops       []Operation `json:"operations"`
	Memo      string      `json:"memo"`      // 备注
	Signature string      `json:"signature"` // hex编码
}

type Operation struct {
	To    string `json:"to"`    // 合约地址
	Value string `json:"value"` //
}

// ----合约相关
type SignedMultisigEvmTx struct {
	Mode     int    `json:"mode"`     // 交易模式，默认为0；0-同步模式 1-全异步 2-半异步；如果tx执行时间较长、网络不稳定、出块慢，建议使用半异步模式。
	Deadline uint64 `json:"deadline"` // 有效截止时间
	GasLimit uint64 `json:"gasLimit"` // gas限额
	GasPrice string `json:"gasPrice"` // gas价格，最低为1
	From     string `json:"from"`     // 多签账户地址
	Nonce    uint64 `json:"nonce"`    // 多签账户nonce
	To       string `json:"to"`       // 交易接受者地址或合约地址
	Value    string `json:"value"`    // 交易金额
	Load     string `json:"load"`     // 合约负载，普通原声币转账时为空
	Memo     string `json:"memo"`     // 备注

	Signature MultiSignature `json:"signature"` // 交易签名
}

type MultiSignature struct {
	PubKey     PubKeyMultisigThreshold `json:"pubKey"`     // 多签公钥
	Signatures []string                `json:"signatures"` // 用户签名列表
}

type PubKeyMultisigThreshold struct {
	K       int      `json:"threshold"` // 多签阈值
	PubKeys []string `json:"pubkeys"`   // 多签用户公钥列表
}
