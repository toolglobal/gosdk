# mondo golang-sdk

[TOC]

## 重要说明
- REST-API SDK访问mondo 的API程序，查询代币转账记录需要API事先配置代币信息生成转账记录
- 执行OLO转账时，接收方账户如果不存在会自动创建，不额外消耗gas，因此OLO转账使用固定gas：21000
- 执行代币转账时，接收方账户如果不存在会自动创建，额外消耗gas数量32000，代币转账gas建议值：100000
- 执行代币转账自动创建的账户没有OLO账户，需要额外创建

## 方法索引
### REST-API SDK
|      方法      |            功能            |  备注  |
| -------------- | -------------------------- | ----- |
| GenKey         | 生成OLO以及基于OLO的代币账户 |  |
| ValidAddress   | 检查是否式合法的OLO地址      |       |
| ValidPublicKey | 检查是否式合法的OLO公钥      |       |
| Exist          | 检查账户是否已在链上存在     |       |
| GetBalance     | 查询OLO账户余额             |       |
| Payment        | OLO转账                    |       |
| Payments       | 批量OLO转账                 |       |
| V3GetLedgers            | 查询所有账本（区块）         |       |
| V3GetLedger             | 根据高度查询账本（区块）     |       |
| V3GetLedgerTransactions | 根据高度查询交易            |       |
| V3GetLedgerPayments     | 根据高度查询转账记录         |       |
| V3GetTransactions       | 代币交易                    |       |
| V3GetTransaction        | 指定交易hash查询交易        |       |
| V3GetPayments           | 查询所有转账记录            |       |
| V3GetTxPayments         | 根据交易hash查询转账记录     |       |
| V3GetAccountPayments    | 根据账户地址查询转账记录     |       |
| ERC20BalanceOf          | 查询代币余额                |       |
| ERC20Pay                | 代币转账                    |       |
| BuildEvmTx                | 生成EVM交易                    |       |
| SendEvmTx                | 发送EVM交易                    |       |
| CheckTx                | 检查交易是否上链并处理成功                    |       |
| CheckTxWithCtx                | 检查交易是否上链并处理成功                    |       |

