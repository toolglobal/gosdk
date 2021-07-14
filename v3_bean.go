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
