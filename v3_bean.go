package gosdk

import (
	"math/big"
	"time"
)

type V3Ledger struct {
	Id         uint64    `db:"id" json:"id"`
	Height     int64     `db:"height" json:"height"`
	BlockHash  string    `db:"blockHash" json:"blockHash"`
	BlockSize  int       `db:"blockSize" json:"blockSize"`
	Validator  string    `db:"validator" json:"validator"`
	TxCount    int64     `db:"txCount" json:"txCount"`
	GasLimit   int64     `db:"gasLimit" json:"gasLimit"`
	GasUsed    int64     `db:"gasUsed" json:"gasUsed"`
	GasPrice   string    `db:"gasPrice" json:"gasPrice"`
	CreatedAt  time.Time `db:"createdAt" json:"createdAt"`
	TotalPrice *big.Int  `db:"-" json:"-"`
}

type V3Transaction struct {
	Id        uint64    `db:"id" json:"id"`
	Hash      string    `db:"hash" json:"hash"`
	Height    int64     `db:"height" json:"height"`
	Typei     int       `db:"typei" json:"typei"`
	Types     string    `db:"types" json:"types"`
	Sender    string    `db:"sender" json:"sender"`
	Nonce     int64     `db:"nonce" json:"nonce"`
	Receiver  string    `db:"receiver" json:"receiver"`
	Value     string    `db:"value" json:"value"`
	GasLimit  int64     `db:"gasLimit" json:"gasLimit"`
	GasUsed   int64     `db:"gasUsed" json:"gasUsed"`
	GasPrice  string    `db:"gasPrice" json:"gasPrice"`
	Memo      string    `db:"memo" json:"memo"`
	Payload   string    `db:"payload" json:"payload"`
	Events    string    `db:"events" json:"events"`
	Codei     uint32    `db:"codei" json:"codei"`
	Codes     string    `db:"codes" json:"codes"`
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
}

type V3Payment struct {
	Id        uint64    `db:"id" json:"id"`
	Hash      string    `db:"hash" json:"hash"`
	Height    int64     `db:"height" json:"height"`
	Idx       int       `db:"idx" json:"idx"`
	Sender    string    `db:"sender" json:"sender"`
	Receiver  string    `db:"receiver" json:"receiver"`
	Symbol    string    `db:"symbol" json:"symbol"`
	Contract  string    `db:"contract" json:"contract"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
}

type BeanValidator struct {
	Address    string `json:"address"`    // 节点公钥
	Nonce      uint64 `json:"nonce"`      // nonce
	Power      uint64 `json:"power"`      // 投票权
	Genesis    bool   `json:"genesis"`    // 是否是创世节点
	RunFor     bool   `json:"runFor"`     // 参与竞选
	Balance    string `json:"balance"`    // 余额（累积收益）
	Profit     string `json:"-"`          // 本轮收益
	Mortgaged  string `json:"mortgaged"`  // 已抵押
	ToMortgage string `json:"toMortgage"` // 待抵押
	Redeem     bool   `json:"redeem"`     // 待赎回
	Voters     int    `json:"voters"`     // 投票人数
	Voted      string `json:"voted"`      // 选民贡献
}

type BeanVoter struct {
	Address    string `json:"address"`    // 选民地址
	Balance    string `json:"balance"`    // 累积收益余额
	Validator  string `json:"validator"`  // 节点地址
	Mortgaged  string `json:"mortgaged"`  // 已抵押
	ToMortgage string `json:"toMortgage"` // 待抵押
	Burned     string `json:"burned"`     // 已燃烧
	Redeem     bool   `json:"redeem"`     // 赎回标志
}

type BeanVots struct {
	Address string    `json:"address"` // 节点公钥
	Voted   string    `json:"voted"`   // 选民贡献
	Voters  []BeanVot `json:"voters"`  // 选民列表
}

type BeanVot struct {
	Address string `json:"address"` // 选民地址
	Amount  string `json:"amount"`  // 投票金额
}
