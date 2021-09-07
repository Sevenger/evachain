# evachain

用 Golang 实现的 blockchain。

## 资料

- [github: awesome-blockchain-cn](https://github.com/chaozh/awesome-blockchain-cn)
- [github: awesome-blockchain](https://github.com/yjjnls/awesome-blockchain)
- [Medium: A blockchain in 200 lines of code](https://medium.com/@lhartikk/a-blockchain-in-200-lines-of-code-963cc1cc0e54)

## 区块链的核心

### 基础架构

- 网络层：p2p、websocket信息传输
- 数据层：数据持久化、哈希函数
- 共识层：POW、POS等
- 应用层：区块链客户端

### abstract

Block、Blockchain

```go
type IBlock interface {
	IsValidBlock(lastBlock IBlock) bool 
	CalculateHash() string             
}

type IBlockchain interface {
	IsValidChain() bool
	Add、GetBlock()
	GenerateBlock() IBlock
	ReplaceChain()
}

type BlockImpl struct {
	Hash、PreHash
	High
	Data: {
	    Tx: {
	        Address: { From, To }
        }
	    Value: { Token, Signature }
    }
}

type BlockchainImpl struct {
	Blocks []IBlock
}
```

Node

```go
type INode interface {
	Connect(addr string)  // 链接节点
	BroadCast(msg string) // 广播消息
    HandleFunc()... // Handle函数
}
```
