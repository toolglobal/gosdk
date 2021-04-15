package gosdk

import "time"

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

type DelegateTx struct {
	CreatedAt uint64 `json:"createdAt"` // 时间戳unixNano
	Sender    string `json:"sender"`    // pubkey
	Nonce     uint64 `json:"nonce"`     //
	OpType    uint8  `json:"opType"`    // 1-抵押选举 2-赎回 3-领取收益
	OpValue   string `json:"opValue"`   // Op对应的值
	Receiver  string `json:"receiver"`  // 选举节点地址
	Signature string `json:"signature"` // hex编码
}

type DPOSPool struct {
	CurrentHeight int64 `json:"currentHeight"` // 当前高度
	LastHeight    int64 `json:"lastHeight"`    // 上一轮DPOS高度
	NextHeight    int64 `json:"nextHeight"`    // 下一轮DPOS高度
}

type DPOSPoolLog struct {
	Id       uint64    `db:"id" json:"id"`             // 数据库自增id
	Height   uint64    `db:"height" json:"height"`     // 去看高度
	Balance  string    `db:"balance" json:"balance"`   // 矿池余额
	Mined    string    `db:"mined" json:"mined"`       // 挖矿量
	Released string    `db:"released" json:"released"` // 实际释放量，实际释放量约等于挖矿量
	Total    string    `db:"total" json:"total"`       // 总股权
	BlockAt  time.Time `db:"blockat" json:"blockAt"`   // 区块时间
}

type DPOSTcnLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 节点抵押量
	Voted     string    `db:"voted" json:"voted"`         // 用户投票量
	Voters    uint64    `db:"voters" json:"voters"`       // 投票用户数
	Profit    string    `db:"profit" json:"profit"`       // 收益
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
}

type DPOSTinLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 用户地址
	Validator string    `db:"validator" json:"validator"` // 用户选举的节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 用户抵押量
	Profit    string    `db:"profit" json:"profit"`       // 收益
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
}

type DPOSRankLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 节点抵押量
	Voted     string    `db:"voted" json:"voted"`         // 用户投票量
	Voters    uint64    `db:"voters" json:"voters"`       // 投票用户数
	Total     string    `db:"total" json:"total"`         // 总股权
	Rank      uint32    `db:"rank" json:"rank"`           // 排名
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
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
