package mempool

import (
	"github.com/Rachelgill00/simulation-for-calvin/config"
	"github.com/Rachelgill00/simulation-for-calvin/message"
)

// Producer结构体包含 一个名为mempool的成员，它是一个指向 MemPool类型的指针。
// 这意味着每个Producer实例 都拥有一个 内存池对象，可以 与之交互 来 管理交易。
type Producer struct {
	mempool *MemPool
}

func NewProducer() *Producer {
	return &Producer{
		mempool: NewMemPool(),
	}
}

// mempool에서 transaction payload 생선
func (pd *Producer) GeneratePayload() []*message.Transaction {
	return pd.mempool.some(config.Configuration.BSize)
}

// 한 새로운 txn mempool에서 추가 (new)
func (pd *Producer) AddTxn(txn *message.Transaction) {
	pd.mempool.addNew(txn) //new
}

// 한 새로운 txn mempool에서 추가 (old)
func (pd *Producer) CollectTxn(txn *message.Transaction) {
	pd.mempool.addOld(txn) //old
}

// return total of txn
func (pd *Producer) TotalReceivedTxNo() int64 {
	return pd.mempool.totalReceived
}
